// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/devldavydov/myfood/internal/storage/ent/migrate"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/devldavydov/myfood/internal/storage/ent/food"
	"github.com/devldavydov/myfood/internal/storage/ent/journal"
	"github.com/devldavydov/myfood/internal/storage/ent/usersettings"
	"github.com/devldavydov/myfood/internal/storage/ent/weight"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// Food is the client for interacting with the Food builders.
	Food *FoodClient
	// Journal is the client for interacting with the Journal builders.
	Journal *JournalClient
	// UserSettings is the client for interacting with the UserSettings builders.
	UserSettings *UserSettingsClient
	// Weight is the client for interacting with the Weight builders.
	Weight *WeightClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	client := &Client{config: newConfig(opts...)}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.Food = NewFoodClient(c.config)
	c.Journal = NewJournalClient(c.config)
	c.UserSettings = NewUserSettingsClient(c.config)
	c.Weight = NewWeightClient(c.config)
}

type (
	// config is the configuration for the client and its builder.
	config struct {
		// driver used for executing database requests.
		driver dialect.Driver
		// debug enable a debug logging.
		debug bool
		// log used for logging on debug mode.
		log func(...any)
		// hooks to execute on mutations.
		hooks *hooks
		// interceptors to execute on queries.
		inters *inters
	}
	// Option function to configure the client.
	Option func(*config)
)

// newConfig creates a new config for the client.
func newConfig(opts ...Option) config {
	cfg := config{log: log.Println, hooks: &hooks{}, inters: &inters{}}
	cfg.options(opts...)
	return cfg
}

// options applies the options on the config object.
func (c *config) options(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
	if c.debug {
		c.driver = dialect.Debug(c.driver, c.log)
	}
}

// Debug enables debug logging on the ent.Driver.
func Debug() Option {
	return func(c *config) {
		c.debug = true
	}
}

// Log sets the logging function for debug mode.
func Log(fn func(...any)) Option {
	return func(c *config) {
		c.log = fn
	}
}

// Driver configures the client driver.
func Driver(driver dialect.Driver) Option {
	return func(c *config) {
		c.driver = driver
	}
}

// Open opens a database/sql.DB specified by the driver name and
// the data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// ErrTxStarted is returned when trying to start a new transaction from a transactional client.
var ErrTxStarted = errors.New("ent: cannot start a transaction within a transaction")

// Tx returns a new transactional client. The provided context
// is used until the transaction is committed or rolled back.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, ErrTxStarted
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = tx
	return &Tx{
		ctx:          ctx,
		config:       cfg,
		Food:         NewFoodClient(cfg),
		Journal:      NewJournalClient(cfg),
		UserSettings: NewUserSettingsClient(cfg),
		Weight:       NewWeightClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, errors.New("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	}).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = &txDriver{tx: tx, drv: c.driver}
	return &Tx{
		ctx:          ctx,
		config:       cfg,
		Food:         NewFoodClient(cfg),
		Journal:      NewJournalClient(cfg),
		UserSettings: NewUserSettingsClient(cfg),
		Weight:       NewWeightClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		Food.
//		Query().
//		Count(ctx)
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := c.config
	cfg.driver = dialect.Debug(c.driver, c.log)
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	c.Food.Use(hooks...)
	c.Journal.Use(hooks...)
	c.UserSettings.Use(hooks...)
	c.Weight.Use(hooks...)
}

// Intercept adds the query interceptors to all the entity clients.
// In order to add interceptors to a specific client, call: `client.Node.Intercept(...)`.
func (c *Client) Intercept(interceptors ...Interceptor) {
	c.Food.Intercept(interceptors...)
	c.Journal.Intercept(interceptors...)
	c.UserSettings.Intercept(interceptors...)
	c.Weight.Intercept(interceptors...)
}

// Mutate implements the ent.Mutator interface.
func (c *Client) Mutate(ctx context.Context, m Mutation) (Value, error) {
	switch m := m.(type) {
	case *FoodMutation:
		return c.Food.mutate(ctx, m)
	case *JournalMutation:
		return c.Journal.mutate(ctx, m)
	case *UserSettingsMutation:
		return c.UserSettings.mutate(ctx, m)
	case *WeightMutation:
		return c.Weight.mutate(ctx, m)
	default:
		return nil, fmt.Errorf("ent: unknown mutation type %T", m)
	}
}

// FoodClient is a client for the Food schema.
type FoodClient struct {
	config
}

// NewFoodClient returns a client for the Food from the given config.
func NewFoodClient(c config) *FoodClient {
	return &FoodClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `food.Hooks(f(g(h())))`.
func (c *FoodClient) Use(hooks ...Hook) {
	c.hooks.Food = append(c.hooks.Food, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `food.Intercept(f(g(h())))`.
func (c *FoodClient) Intercept(interceptors ...Interceptor) {
	c.inters.Food = append(c.inters.Food, interceptors...)
}

// Create returns a builder for creating a Food entity.
func (c *FoodClient) Create() *FoodCreate {
	mutation := newFoodMutation(c.config, OpCreate)
	return &FoodCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Food entities.
func (c *FoodClient) CreateBulk(builders ...*FoodCreate) *FoodCreateBulk {
	return &FoodCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *FoodClient) MapCreateBulk(slice any, setFunc func(*FoodCreate, int)) *FoodCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &FoodCreateBulk{err: fmt.Errorf("calling to FoodClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*FoodCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &FoodCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Food.
func (c *FoodClient) Update() *FoodUpdate {
	mutation := newFoodMutation(c.config, OpUpdate)
	return &FoodUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *FoodClient) UpdateOne(f *Food) *FoodUpdateOne {
	mutation := newFoodMutation(c.config, OpUpdateOne, withFood(f))
	return &FoodUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *FoodClient) UpdateOneID(id int) *FoodUpdateOne {
	mutation := newFoodMutation(c.config, OpUpdateOne, withFoodID(id))
	return &FoodUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Food.
func (c *FoodClient) Delete() *FoodDelete {
	mutation := newFoodMutation(c.config, OpDelete)
	return &FoodDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *FoodClient) DeleteOne(f *Food) *FoodDeleteOne {
	return c.DeleteOneID(f.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *FoodClient) DeleteOneID(id int) *FoodDeleteOne {
	builder := c.Delete().Where(food.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &FoodDeleteOne{builder}
}

// Query returns a query builder for Food.
func (c *FoodClient) Query() *FoodQuery {
	return &FoodQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeFood},
		inters: c.Interceptors(),
	}
}

// Get returns a Food entity by its id.
func (c *FoodClient) Get(ctx context.Context, id int) (*Food, error) {
	return c.Query().Where(food.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *FoodClient) GetX(ctx context.Context, id int) *Food {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryJournals queries the journals edge of a Food.
func (c *FoodClient) QueryJournals(f *Food) *JournalQuery {
	query := (&JournalClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := f.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(food.Table, food.FieldID, id),
			sqlgraph.To(journal.Table, journal.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, food.JournalsTable, food.JournalsColumn),
		)
		fromV = sqlgraph.Neighbors(f.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *FoodClient) Hooks() []Hook {
	return c.hooks.Food
}

// Interceptors returns the client interceptors.
func (c *FoodClient) Interceptors() []Interceptor {
	return c.inters.Food
}

func (c *FoodClient) mutate(ctx context.Context, m *FoodMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&FoodCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&FoodUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&FoodUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&FoodDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Food mutation op: %q", m.Op())
	}
}

// JournalClient is a client for the Journal schema.
type JournalClient struct {
	config
}

// NewJournalClient returns a client for the Journal from the given config.
func NewJournalClient(c config) *JournalClient {
	return &JournalClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `journal.Hooks(f(g(h())))`.
func (c *JournalClient) Use(hooks ...Hook) {
	c.hooks.Journal = append(c.hooks.Journal, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `journal.Intercept(f(g(h())))`.
func (c *JournalClient) Intercept(interceptors ...Interceptor) {
	c.inters.Journal = append(c.inters.Journal, interceptors...)
}

// Create returns a builder for creating a Journal entity.
func (c *JournalClient) Create() *JournalCreate {
	mutation := newJournalMutation(c.config, OpCreate)
	return &JournalCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Journal entities.
func (c *JournalClient) CreateBulk(builders ...*JournalCreate) *JournalCreateBulk {
	return &JournalCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *JournalClient) MapCreateBulk(slice any, setFunc func(*JournalCreate, int)) *JournalCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &JournalCreateBulk{err: fmt.Errorf("calling to JournalClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*JournalCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &JournalCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Journal.
func (c *JournalClient) Update() *JournalUpdate {
	mutation := newJournalMutation(c.config, OpUpdate)
	return &JournalUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *JournalClient) UpdateOne(j *Journal) *JournalUpdateOne {
	mutation := newJournalMutation(c.config, OpUpdateOne, withJournal(j))
	return &JournalUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *JournalClient) UpdateOneID(id int) *JournalUpdateOne {
	mutation := newJournalMutation(c.config, OpUpdateOne, withJournalID(id))
	return &JournalUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Journal.
func (c *JournalClient) Delete() *JournalDelete {
	mutation := newJournalMutation(c.config, OpDelete)
	return &JournalDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *JournalClient) DeleteOne(j *Journal) *JournalDeleteOne {
	return c.DeleteOneID(j.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *JournalClient) DeleteOneID(id int) *JournalDeleteOne {
	builder := c.Delete().Where(journal.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &JournalDeleteOne{builder}
}

// Query returns a query builder for Journal.
func (c *JournalClient) Query() *JournalQuery {
	return &JournalQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeJournal},
		inters: c.Interceptors(),
	}
}

// Get returns a Journal entity by its id.
func (c *JournalClient) Get(ctx context.Context, id int) (*Journal, error) {
	return c.Query().Where(journal.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *JournalClient) GetX(ctx context.Context, id int) *Journal {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryFood queries the food edge of a Journal.
func (c *JournalClient) QueryFood(j *Journal) *FoodQuery {
	query := (&FoodClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := j.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(journal.Table, journal.FieldID, id),
			sqlgraph.To(food.Table, food.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, journal.FoodTable, journal.FoodColumn),
		)
		fromV = sqlgraph.Neighbors(j.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *JournalClient) Hooks() []Hook {
	return c.hooks.Journal
}

// Interceptors returns the client interceptors.
func (c *JournalClient) Interceptors() []Interceptor {
	return c.inters.Journal
}

func (c *JournalClient) mutate(ctx context.Context, m *JournalMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&JournalCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&JournalUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&JournalUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&JournalDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Journal mutation op: %q", m.Op())
	}
}

// UserSettingsClient is a client for the UserSettings schema.
type UserSettingsClient struct {
	config
}

// NewUserSettingsClient returns a client for the UserSettings from the given config.
func NewUserSettingsClient(c config) *UserSettingsClient {
	return &UserSettingsClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `usersettings.Hooks(f(g(h())))`.
func (c *UserSettingsClient) Use(hooks ...Hook) {
	c.hooks.UserSettings = append(c.hooks.UserSettings, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `usersettings.Intercept(f(g(h())))`.
func (c *UserSettingsClient) Intercept(interceptors ...Interceptor) {
	c.inters.UserSettings = append(c.inters.UserSettings, interceptors...)
}

// Create returns a builder for creating a UserSettings entity.
func (c *UserSettingsClient) Create() *UserSettingsCreate {
	mutation := newUserSettingsMutation(c.config, OpCreate)
	return &UserSettingsCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of UserSettings entities.
func (c *UserSettingsClient) CreateBulk(builders ...*UserSettingsCreate) *UserSettingsCreateBulk {
	return &UserSettingsCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *UserSettingsClient) MapCreateBulk(slice any, setFunc func(*UserSettingsCreate, int)) *UserSettingsCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &UserSettingsCreateBulk{err: fmt.Errorf("calling to UserSettingsClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*UserSettingsCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &UserSettingsCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for UserSettings.
func (c *UserSettingsClient) Update() *UserSettingsUpdate {
	mutation := newUserSettingsMutation(c.config, OpUpdate)
	return &UserSettingsUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *UserSettingsClient) UpdateOne(us *UserSettings) *UserSettingsUpdateOne {
	mutation := newUserSettingsMutation(c.config, OpUpdateOne, withUserSettings(us))
	return &UserSettingsUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *UserSettingsClient) UpdateOneID(id int) *UserSettingsUpdateOne {
	mutation := newUserSettingsMutation(c.config, OpUpdateOne, withUserSettingsID(id))
	return &UserSettingsUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for UserSettings.
func (c *UserSettingsClient) Delete() *UserSettingsDelete {
	mutation := newUserSettingsMutation(c.config, OpDelete)
	return &UserSettingsDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *UserSettingsClient) DeleteOne(us *UserSettings) *UserSettingsDeleteOne {
	return c.DeleteOneID(us.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *UserSettingsClient) DeleteOneID(id int) *UserSettingsDeleteOne {
	builder := c.Delete().Where(usersettings.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &UserSettingsDeleteOne{builder}
}

// Query returns a query builder for UserSettings.
func (c *UserSettingsClient) Query() *UserSettingsQuery {
	return &UserSettingsQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeUserSettings},
		inters: c.Interceptors(),
	}
}

// Get returns a UserSettings entity by its id.
func (c *UserSettingsClient) Get(ctx context.Context, id int) (*UserSettings, error) {
	return c.Query().Where(usersettings.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *UserSettingsClient) GetX(ctx context.Context, id int) *UserSettings {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *UserSettingsClient) Hooks() []Hook {
	return c.hooks.UserSettings
}

// Interceptors returns the client interceptors.
func (c *UserSettingsClient) Interceptors() []Interceptor {
	return c.inters.UserSettings
}

func (c *UserSettingsClient) mutate(ctx context.Context, m *UserSettingsMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&UserSettingsCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&UserSettingsUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&UserSettingsUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&UserSettingsDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown UserSettings mutation op: %q", m.Op())
	}
}

// WeightClient is a client for the Weight schema.
type WeightClient struct {
	config
}

// NewWeightClient returns a client for the Weight from the given config.
func NewWeightClient(c config) *WeightClient {
	return &WeightClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `weight.Hooks(f(g(h())))`.
func (c *WeightClient) Use(hooks ...Hook) {
	c.hooks.Weight = append(c.hooks.Weight, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `weight.Intercept(f(g(h())))`.
func (c *WeightClient) Intercept(interceptors ...Interceptor) {
	c.inters.Weight = append(c.inters.Weight, interceptors...)
}

// Create returns a builder for creating a Weight entity.
func (c *WeightClient) Create() *WeightCreate {
	mutation := newWeightMutation(c.config, OpCreate)
	return &WeightCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Weight entities.
func (c *WeightClient) CreateBulk(builders ...*WeightCreate) *WeightCreateBulk {
	return &WeightCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *WeightClient) MapCreateBulk(slice any, setFunc func(*WeightCreate, int)) *WeightCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &WeightCreateBulk{err: fmt.Errorf("calling to WeightClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*WeightCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &WeightCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Weight.
func (c *WeightClient) Update() *WeightUpdate {
	mutation := newWeightMutation(c.config, OpUpdate)
	return &WeightUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *WeightClient) UpdateOne(w *Weight) *WeightUpdateOne {
	mutation := newWeightMutation(c.config, OpUpdateOne, withWeight(w))
	return &WeightUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *WeightClient) UpdateOneID(id int) *WeightUpdateOne {
	mutation := newWeightMutation(c.config, OpUpdateOne, withWeightID(id))
	return &WeightUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Weight.
func (c *WeightClient) Delete() *WeightDelete {
	mutation := newWeightMutation(c.config, OpDelete)
	return &WeightDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *WeightClient) DeleteOne(w *Weight) *WeightDeleteOne {
	return c.DeleteOneID(w.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *WeightClient) DeleteOneID(id int) *WeightDeleteOne {
	builder := c.Delete().Where(weight.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &WeightDeleteOne{builder}
}

// Query returns a query builder for Weight.
func (c *WeightClient) Query() *WeightQuery {
	return &WeightQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeWeight},
		inters: c.Interceptors(),
	}
}

// Get returns a Weight entity by its id.
func (c *WeightClient) Get(ctx context.Context, id int) (*Weight, error) {
	return c.Query().Where(weight.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *WeightClient) GetX(ctx context.Context, id int) *Weight {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// Hooks returns the client hooks.
func (c *WeightClient) Hooks() []Hook {
	return c.hooks.Weight
}

// Interceptors returns the client interceptors.
func (c *WeightClient) Interceptors() []Interceptor {
	return c.inters.Weight
}

func (c *WeightClient) mutate(ctx context.Context, m *WeightMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&WeightCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&WeightUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&WeightUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&WeightDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Weight mutation op: %q", m.Op())
	}
}

// hooks and interceptors per client, for fast access.
type (
	hooks struct {
		Food, Journal, UserSettings, Weight []ent.Hook
	}
	inters struct {
		Food, Journal, UserSettings, Weight []ent.Interceptor
	}
)
