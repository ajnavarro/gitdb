package git

import (
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/ajnavarro/gitdb/model"
	"github.com/ajnavarro/gitdb/ops"
	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/filemode"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

const dataBlobName = "ops"

type Repository struct {
	repo *gogit.Repository
}

func NewRepository(r *gogit.Repository) *Repository {
	return &Repository{r}
}

func (r *Repository) UpdateRow(rowID, dbName, tableName string, data []byte, author *model.Author) error {
	blob, err := r.blob(data)
	if err != nil {
		return err
	}

	ops := r.treeEntry(dataBlobName, blob.Hash())

	root, err := r.root(ops)
	if err != nil {
		return err
	}

	old, err := r.getReference(dbName, tableName, rowID)

	oldHash := old.Hash()

	commit, err := r.commit(author, root.Hash(), &oldHash)
	if err != nil {
		return err
	}

	if err := r.saveObjects(blob, root, commit); err != nil {
		return err
	}

	return r.updateReference(old, commit.Hash())
}

func (r *Repository) NewRow(dbName, tableName string, data []byte, author *model.Author) (string, error) {
	blob, err := r.blob(data)
	if err != nil {
		return "", err
	}

	ops := r.treeEntry(dataBlobName, blob.Hash())

	root, err := r.root(ops)
	if err != nil {
		return "", err
	}

	commit, err := r.commit(author, root.Hash(), nil)
	if err != nil {
		return "", err
	}

	if err := r.saveObjects(blob, root, commit); err != nil {
		return "", err
	}

	if err := r.createReference(dbName, tableName, commit.Hash()); err != nil {
		return "", err
	}

	return commit.Hash().String(), nil
}

func (r *Repository) GetOperationBlocks(rowID, dbName, tableName string) (ops.OperationBlockIter, error) {
	ref, err := r.getReference(dbName, tableName, rowID)
	if err != nil {
		return nil, err
	}

	return &DataRowIter{r.repo, ref.Hash()}, nil
}

type DataRowIter struct {
	repo *gogit.Repository
	next plumbing.Hash
}

func (i *DataRowIter) Next() (*model.OperationBlock, error) {
	if i.next == plumbing.ZeroHash {
		return nil, io.EOF
	}

	c, err := i.repo.CommitObject(i.next)
	if err != nil {
		return nil, err
	}

	f, err := c.File(dataBlobName)
	if err != nil {
		return nil, err
	}

	fr, err := f.Blob.Reader()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(fr)
	if err != nil {
		return nil, err
	}

	parents := c.ParentHashes
	if len(parents) > 1 {
		return nil, fmt.Errorf("commit history with several branches not supported")
	}

	if len(parents) == 0 {
		i.next = plumbing.ZeroHash
	} else {
		i.next = parents[0]
	}

	return model.UnmarshalOperationBlock(data)
}

func (i *DataRowIter) Close() error {
	i.repo = nil
	i.next = plumbing.ZeroHash

	return nil
}

func (r *Repository) getReference(dbName, tableName, rowID string) (*plumbing.Reference, error) {
	refName := plumbing.ReferenceName(
		model.RowReference(
			dbName,
			tableName,
			rowID,
		),
	)

	return r.repo.Reference(refName, true)
}

func (r *Repository) updateReference(old *plumbing.Reference, commit plumbing.Hash) error {
	new := plumbing.NewHashReference(old.Name(), commit)

	return r.repo.Storer.SetReference(new)
}
func (r *Repository) createReference(dbName, tableName string, commitHash plumbing.Hash) error {
	refName := plumbing.ReferenceName(
		model.RowReference(
			dbName,
			tableName,
			commitHash.String(),
		),
	)

	ref := plumbing.NewHashReference(
		refName,
		commitHash,
	)

	return r.repo.Storer.SetReference(ref)
}

func (r *Repository) saveObjects(objects ...plumbing.EncodedObject) error {

	for _, o := range objects {
		if _, err := r.repo.Storer.SetEncodedObject(o); err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) commit(author *model.Author, rootHash plumbing.Hash, parent *plumbing.Hash) (plumbing.EncodedObject, error) {
	sig := object.Signature{
		Email: author.Email,
		Name:  author.Name,
		// TODO we cannot trust timestamps for nodes.
		When: time.Now(),
	}

	commit := &object.Commit{
		Author:    sig,
		Committer: sig,
		TreeHash:  rootHash,
	}

	if parent != nil {
		commit.ParentHashes = []plumbing.Hash{*parent}
	}

	commitObj := &plumbing.MemoryObject{}
	if err := commit.Encode(commitObj); err != nil {
		return nil, err
	}

	return commitObj, nil
}

func (r *Repository) root(tes ...object.TreeEntry) (plumbing.EncodedObject, error) {
	// create root tree
	root := &object.Tree{
		Entries: tes,
	}

	treeEntryObj := &plumbing.MemoryObject{}
	if err := root.Encode(treeEntryObj); err != nil {
		return nil, err
	}

	return treeEntryObj, nil
}

func (r *Repository) treeEntry(name string, blobHash plumbing.Hash) object.TreeEntry {
	// create tree entry
	teOps := object.TreeEntry{
		Name: name,
		Mode: filemode.Regular,
		Hash: blobHash,
	}

	return teOps
}

func (r *Repository) blob(data []byte) (plumbing.EncodedObject, error) {
	// create blob
	blobObj := &plumbing.MemoryObject{}
	blobObj.SetType(plumbing.BlobObject)
	if _, err := blobObj.Write(data); err != nil {
		return nil, err
	}

	b := &object.Blob{}
	if err := b.Decode(blobObj); err != nil {
		return nil, err
	}

	return blobObj, nil
}
