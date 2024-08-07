// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/devldavydov/myfood/internal/storage/ent/bundle"
	"github.com/devldavydov/myfood/internal/storage/ent/predicate"
)

// BundleQuery is the builder for querying Bundle entities.
type BundleQuery struct {
	config
	ctx        *QueryContext
	order      []bundle.OrderOption
	inters     []Interceptor
	predicates []predicate.Bundle
	modifiers  []func(*sql.Selector)
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the BundleQuery builder.
func (bq *BundleQuery) Where(ps ...predicate.Bundle) *BundleQuery {
	bq.predicates = append(bq.predicates, ps...)
	return bq
}

// Limit the number of records to be returned by this query.
func (bq *BundleQuery) Limit(limit int) *BundleQuery {
	bq.ctx.Limit = &limit
	return bq
}

// Offset to start from.
func (bq *BundleQuery) Offset(offset int) *BundleQuery {
	bq.ctx.Offset = &offset
	return bq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (bq *BundleQuery) Unique(unique bool) *BundleQuery {
	bq.ctx.Unique = &unique
	return bq
}

// Order specifies how the records should be ordered.
func (bq *BundleQuery) Order(o ...bundle.OrderOption) *BundleQuery {
	bq.order = append(bq.order, o...)
	return bq
}

// First returns the first Bundle entity from the query.
// Returns a *NotFoundError when no Bundle was found.
func (bq *BundleQuery) First(ctx context.Context) (*Bundle, error) {
	nodes, err := bq.Limit(1).All(setContextOp(ctx, bq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{bundle.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (bq *BundleQuery) FirstX(ctx context.Context) *Bundle {
	node, err := bq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Bundle ID from the query.
// Returns a *NotFoundError when no Bundle ID was found.
func (bq *BundleQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = bq.Limit(1).IDs(setContextOp(ctx, bq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{bundle.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (bq *BundleQuery) FirstIDX(ctx context.Context) int {
	id, err := bq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single Bundle entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one Bundle entity is found.
// Returns a *NotFoundError when no Bundle entities are found.
func (bq *BundleQuery) Only(ctx context.Context) (*Bundle, error) {
	nodes, err := bq.Limit(2).All(setContextOp(ctx, bq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{bundle.Label}
	default:
		return nil, &NotSingularError{bundle.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (bq *BundleQuery) OnlyX(ctx context.Context) *Bundle {
	node, err := bq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only Bundle ID in the query.
// Returns a *NotSingularError when more than one Bundle ID is found.
// Returns a *NotFoundError when no entities are found.
func (bq *BundleQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = bq.Limit(2).IDs(setContextOp(ctx, bq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{bundle.Label}
	default:
		err = &NotSingularError{bundle.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (bq *BundleQuery) OnlyIDX(ctx context.Context) int {
	id, err := bq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Bundles.
func (bq *BundleQuery) All(ctx context.Context) ([]*Bundle, error) {
	ctx = setContextOp(ctx, bq.ctx, "All")
	if err := bq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*Bundle, *BundleQuery]()
	return withInterceptors[[]*Bundle](ctx, bq, qr, bq.inters)
}

// AllX is like All, but panics if an error occurs.
func (bq *BundleQuery) AllX(ctx context.Context) []*Bundle {
	nodes, err := bq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Bundle IDs.
func (bq *BundleQuery) IDs(ctx context.Context) (ids []int, err error) {
	if bq.ctx.Unique == nil && bq.path != nil {
		bq.Unique(true)
	}
	ctx = setContextOp(ctx, bq.ctx, "IDs")
	if err = bq.Select(bundle.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (bq *BundleQuery) IDsX(ctx context.Context) []int {
	ids, err := bq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (bq *BundleQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, bq.ctx, "Count")
	if err := bq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, bq, querierCount[*BundleQuery](), bq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (bq *BundleQuery) CountX(ctx context.Context) int {
	count, err := bq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (bq *BundleQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, bq.ctx, "Exist")
	switch _, err := bq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (bq *BundleQuery) ExistX(ctx context.Context) bool {
	exist, err := bq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the BundleQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (bq *BundleQuery) Clone() *BundleQuery {
	if bq == nil {
		return nil
	}
	return &BundleQuery{
		config:     bq.config,
		ctx:        bq.ctx.Clone(),
		order:      append([]bundle.OrderOption{}, bq.order...),
		inters:     append([]Interceptor{}, bq.inters...),
		predicates: append([]predicate.Bundle{}, bq.predicates...),
		// clone intermediate query.
		sql:  bq.sql.Clone(),
		path: bq.path,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Userid int64 `json:"userid,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Bundle.Query().
//		GroupBy(bundle.FieldUserid).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (bq *BundleQuery) GroupBy(field string, fields ...string) *BundleGroupBy {
	bq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &BundleGroupBy{build: bq}
	grbuild.flds = &bq.ctx.Fields
	grbuild.label = bundle.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Userid int64 `json:"userid,omitempty"`
//	}
//
//	client.Bundle.Query().
//		Select(bundle.FieldUserid).
//		Scan(ctx, &v)
func (bq *BundleQuery) Select(fields ...string) *BundleSelect {
	bq.ctx.Fields = append(bq.ctx.Fields, fields...)
	sbuild := &BundleSelect{BundleQuery: bq}
	sbuild.label = bundle.Label
	sbuild.flds, sbuild.scan = &bq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a BundleSelect configured with the given aggregations.
func (bq *BundleQuery) Aggregate(fns ...AggregateFunc) *BundleSelect {
	return bq.Select().Aggregate(fns...)
}

func (bq *BundleQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range bq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, bq); err != nil {
				return err
			}
		}
	}
	for _, f := range bq.ctx.Fields {
		if !bundle.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if bq.path != nil {
		prev, err := bq.path(ctx)
		if err != nil {
			return err
		}
		bq.sql = prev
	}
	return nil
}

func (bq *BundleQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*Bundle, error) {
	var (
		nodes = []*Bundle{}
		_spec = bq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*Bundle).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &Bundle{config: bq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	if len(bq.modifiers) > 0 {
		_spec.Modifiers = bq.modifiers
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, bq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (bq *BundleQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := bq.querySpec()
	if len(bq.modifiers) > 0 {
		_spec.Modifiers = bq.modifiers
	}
	_spec.Node.Columns = bq.ctx.Fields
	if len(bq.ctx.Fields) > 0 {
		_spec.Unique = bq.ctx.Unique != nil && *bq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, bq.driver, _spec)
}

func (bq *BundleQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(bundle.Table, bundle.Columns, sqlgraph.NewFieldSpec(bundle.FieldID, field.TypeInt))
	_spec.From = bq.sql
	if unique := bq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if bq.path != nil {
		_spec.Unique = true
	}
	if fields := bq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, bundle.FieldID)
		for i := range fields {
			if fields[i] != bundle.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := bq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := bq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := bq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := bq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (bq *BundleQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(bq.driver.Dialect())
	t1 := builder.Table(bundle.Table)
	columns := bq.ctx.Fields
	if len(columns) == 0 {
		columns = bundle.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if bq.sql != nil {
		selector = bq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if bq.ctx.Unique != nil && *bq.ctx.Unique {
		selector.Distinct()
	}
	for _, m := range bq.modifiers {
		m(selector)
	}
	for _, p := range bq.predicates {
		p(selector)
	}
	for _, p := range bq.order {
		p(selector)
	}
	if offset := bq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := bq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// Modify adds a query modifier for attaching custom logic to queries.
func (bq *BundleQuery) Modify(modifiers ...func(s *sql.Selector)) *BundleSelect {
	bq.modifiers = append(bq.modifiers, modifiers...)
	return bq.Select()
}

// BundleGroupBy is the group-by builder for Bundle entities.
type BundleGroupBy struct {
	selector
	build *BundleQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (bgb *BundleGroupBy) Aggregate(fns ...AggregateFunc) *BundleGroupBy {
	bgb.fns = append(bgb.fns, fns...)
	return bgb
}

// Scan applies the selector query and scans the result into the given value.
func (bgb *BundleGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, bgb.build.ctx, "GroupBy")
	if err := bgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*BundleQuery, *BundleGroupBy](ctx, bgb.build, bgb, bgb.build.inters, v)
}

func (bgb *BundleGroupBy) sqlScan(ctx context.Context, root *BundleQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(bgb.fns))
	for _, fn := range bgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*bgb.flds)+len(bgb.fns))
		for _, f := range *bgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*bgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := bgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// BundleSelect is the builder for selecting fields of Bundle entities.
type BundleSelect struct {
	*BundleQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (bs *BundleSelect) Aggregate(fns ...AggregateFunc) *BundleSelect {
	bs.fns = append(bs.fns, fns...)
	return bs
}

// Scan applies the selector query and scans the result into the given value.
func (bs *BundleSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, bs.ctx, "Select")
	if err := bs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*BundleQuery, *BundleSelect](ctx, bs.BundleQuery, bs, bs.inters, v)
}

func (bs *BundleSelect) sqlScan(ctx context.Context, root *BundleQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(bs.fns))
	for _, fn := range bs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*bs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := bs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// Modify adds a query modifier for attaching custom logic to queries.
func (bs *BundleSelect) Modify(modifiers ...func(s *sql.Selector)) *BundleSelect {
	bs.modifiers = append(bs.modifiers, modifiers...)
	return bs
}
