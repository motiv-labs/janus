package godog

import (
	"fmt"
	"sync"

	"github.com/cucumber/messages-go/v10"
	"github.com/hashicorp/go-memdb"
)

const (
	writeMode bool = true
	readMode  bool = false

	tableFeature         string = "feature"
	tableFeatureIndexURI string = "id"

	tablePickle         string = "pickle"
	tablePickleIndexID  string = "id"
	tablePickleIndexURI string = "uri"

	tablePickleStep        string = "pickle_step"
	tablePickleStepIndexID string = "id"

	tablePickleResult              string = "pickle_result"
	tablePickleResultIndexPickleID string = "id"

	tablePickleStepResult                  string = "pickle_step_result"
	tablePickleStepResultIndexPickleStepID string = "id"
	tablePickleStepResultIndexPickleID     string = "pickle_id"
	tablePickleStepResultIndexStatus       string = "status"
)

type storage struct {
	db *memdb.MemDB

	testRunStarted testRunStarted
	lock           *sync.Mutex
}

func newStorage() *storage {
	// Create the DB schema
	schema := memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			tableFeature: {
				Name: tableFeature,
				Indexes: map[string]*memdb.IndexSchema{
					tableFeatureIndexURI: {
						Name:    tableFeatureIndexURI,
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Uri"},
					},
				},
			},
			tablePickle: {
				Name: tablePickle,
				Indexes: map[string]*memdb.IndexSchema{
					tablePickleIndexID: {
						Name:    tablePickleIndexID,
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Id"},
					},
					tablePickleIndexURI: {
						Name:    tablePickleIndexURI,
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Uri"},
					},
				},
			},
			tablePickleStep: {
				Name: tablePickleStep,
				Indexes: map[string]*memdb.IndexSchema{
					tablePickleStepIndexID: {
						Name:    tablePickleStepIndexID,
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Id"},
					},
				},
			},
			tablePickleResult: {
				Name: tablePickleResult,
				Indexes: map[string]*memdb.IndexSchema{
					tablePickleResultIndexPickleID: {
						Name:    tablePickleResultIndexPickleID,
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "PickleID"},
					},
				},
			},
			tablePickleStepResult: {
				Name: tablePickleStepResult,
				Indexes: map[string]*memdb.IndexSchema{
					tablePickleStepResultIndexPickleStepID: {
						Name:    tablePickleStepResultIndexPickleStepID,
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "PickleStepID"},
					},
					tablePickleStepResultIndexPickleID: {
						Name:    tablePickleStepResultIndexPickleID,
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "PickleID"},
					},
					tablePickleStepResultIndexStatus: {
						Name:    tablePickleStepResultIndexStatus,
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: "Status"},
					},
				},
			},
		},
	}

	db, err := memdb.NewMemDB(&schema)
	if err != nil {
		panic(err)
	}

	return &storage{db: db, lock: new(sync.Mutex)}
}

func (s *storage) mustInsertPickle(p *messages.Pickle) {
	txn := s.db.Txn(writeMode)

	if err := txn.Insert(tablePickle, p); err != nil {
		panic(err)
	}

	for _, step := range p.Steps {
		if err := txn.Insert(tablePickleStep, step); err != nil {
			panic(err)
		}
	}

	txn.Commit()
}

func (s *storage) mustGetPickle(id string) *messages.Pickle {
	v := s.mustFirst(tablePickle, tablePickleIndexID, id)
	return v.(*messages.Pickle)
}

func (s *storage) mustGetPickles(uri string) (ps []*messages.Pickle) {
	it := s.mustGet(tablePickle, tablePickleIndexURI, uri)
	for v := it.Next(); v != nil; v = it.Next() {
		ps = append(ps, v.(*messages.Pickle))
	}

	return
}

func (s *storage) mustGetPickleStep(id string) *messages.Pickle_PickleStep {
	v := s.mustFirst(tablePickleStep, tablePickleStepIndexID, id)
	return v.(*messages.Pickle_PickleStep)
}

func (s *storage) mustInsertTestRunStarted(trs testRunStarted) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.testRunStarted = trs
}

func (s *storage) mustGetTestRunStarted() testRunStarted {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.testRunStarted
}

func (s *storage) mustInsertPickleResult(pr pickleResult) {
	s.mustInsert(tablePickleResult, pr)
}

func (s *storage) mustInsertPickleStepResult(psr pickleStepResult) {
	s.mustInsert(tablePickleStepResult, psr)
}

func (s *storage) mustGetPickleResult(id string) pickleResult {
	v := s.mustFirst(tablePickleResult, tablePickleResultIndexPickleID, id)
	return v.(pickleResult)
}

func (s *storage) mustGetPickleResults() (prs []pickleResult) {
	it := s.mustGet(tablePickleResult, tablePickleResultIndexPickleID)
	for v := it.Next(); v != nil; v = it.Next() {
		prs = append(prs, v.(pickleResult))
	}

	return prs
}

func (s *storage) mustGetPickleStepResult(id string) pickleStepResult {
	v := s.mustFirst(tablePickleStepResult, tablePickleStepResultIndexPickleStepID, id)
	return v.(pickleStepResult)
}

func (s *storage) mustGetPickleStepResultsByPickleID(pickleID string) (psrs []pickleStepResult) {
	it := s.mustGet(tablePickleStepResult, tablePickleStepResultIndexPickleID, pickleID)
	for v := it.Next(); v != nil; v = it.Next() {
		psrs = append(psrs, v.(pickleStepResult))
	}

	return psrs
}

func (s *storage) mustGetPickleStepResultsByStatus(status stepResultStatus) (psrs []pickleStepResult) {
	it := s.mustGet(tablePickleStepResult, tablePickleStepResultIndexStatus, status)
	for v := it.Next(); v != nil; v = it.Next() {
		psrs = append(psrs, v.(pickleStepResult))
	}

	return psrs
}

func (s *storage) mustInsertFeature(f *feature) {
	s.mustInsert(tableFeature, f)
}

func (s *storage) mustGetFeature(uri string) *feature {
	v := s.mustFirst(tableFeature, tableFeatureIndexURI, uri)
	return v.(*feature)
}

func (s *storage) mustGetFeatures() (fs []*feature) {
	it := s.mustGet(tableFeature, tableFeatureIndexURI)
	for v := it.Next(); v != nil; v = it.Next() {
		fs = append(fs, v.(*feature))
	}

	return
}

func (s *storage) mustInsert(table string, obj interface{}) {
	txn := s.db.Txn(writeMode)

	if err := txn.Insert(table, obj); err != nil {
		panic(err)
	}

	txn.Commit()
}

func (s *storage) mustFirst(table, index string, args ...interface{}) interface{} {
	txn := s.db.Txn(readMode)
	defer txn.Abort()

	v, err := txn.First(table, index, args...)
	if err != nil {
		panic(err)
	} else if v == nil {
		err = fmt.Errorf("Couldn't find index: %q in table: %q with args: %+v", index, table, args)
		panic(err)
	}

	return v
}

func (s *storage) mustGet(table, index string, args ...interface{}) memdb.ResultIterator {
	txn := s.db.Txn(readMode)
	defer txn.Abort()

	it, err := txn.Get(table, index, args...)
	if err != nil {
		panic(err)
	}

	return it
}
