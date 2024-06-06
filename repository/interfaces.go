// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	CreateEstate(ctx context.Context, input CreateEstateInput) (err error)
	GetEstateById(ctx context.Context, input GetEstateByIdInput) (output GetEstateByIdOutput, err error)
	CountCoordinateTree(ctx context.Context, input CountCoordinateTreeInput) (output CountCoordinateTreeOutput, err error)
	GetPrevNextTree(ctx context.Context, input GetPrevNextTreeInput) (output GetPrevNextTreeOutput, err error)
	CreateTree(ctx context.Context, input CreateTreeInput) (err error)
	GetHeightEstateTrees(ctx context.Context, input GetHeightEstateTreesInput) (output GetHeightEstateTreesOutput, err error)
	StoreMedianEstate(ctx context.Context, input StoreMedianEstateInput) (err error)
	GetAllEstateTrees(ctx context.Context, input GetAllEstateTreesInput) (output GetAllEstateTreesOutput, err error)
}
