package models

// IModel defines the common interface that all database models must implement.
// This interface is used by GenericRepository to provide standard CRUD operations.
type IModel interface {
	// TableName returns the database table name for the model.
	TableName() string
}

// IModelWithID extends IModel for models that have a numeric ID field.
// This is used by GenericRepository for operations that need to access the record's ID.
type IModelWithID interface {
	IModel

	// GetID returns the primary key ID of the record.
	// This method allows GenericRepository to access the ID without knowing the concrete type.
	GetID() uint
}

// IModelWithCenterID extends IModelWithID for models that belong to a center.
// This enables center-scope operations for multi-tenant data isolation.
type IModelWithCenterID interface {
	IModelWithID

	// GetCenterID returns the center_id of the record.
	// Used for multi-tenant query scoping.
	GetCenterID() uint
}
