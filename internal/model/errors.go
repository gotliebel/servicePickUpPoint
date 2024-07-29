package model

import (
	"errors"
)

var (
	ErrOrderIdMissing            = errors.New("order id is missing")
	ErrClientIdMissing           = errors.New("client id is missing")
	ErrIdListMissing             = errors.New("id list is missing")
	ErrDateMissing               = errors.New("date is missing")
	ErrDateUnparsableTemplate    = "cannot parse date. Try using help to find correct format. Occuring error: %s"
	ErrPageMissing               = errors.New("page number is missing")
	ErrStorageTimeIsIncorrect    = errors.New("storage time of unit is already out, before we accept it from courier")
	ErrTooSoonToTakeBackTemplate = "cannot return this order: end of storage time is %s"
	ErrTimeHasExpired            = errors.New("storage time of some order has already expired")
	ErrTimeToReturnHasExpired    = errors.New("this order cannot be returned already")
	ErrOrderAlreadyExists        = errors.New("this order is already in storage")
	ErrOrderAlreadyTaken         = errors.New("this order is already taken by client")
	ErrDoesntExists              = errors.New("cannot find order with this id to handle it to you")
	ErrOrdersByDifferentClients  = errors.New("orders from this list come from different clients")
	ErrOrderMadeByOtherClient    = errors.New("this order was made by another client, we can't return this to you")
	ErrSomeOrdersWereNotFound    = errors.New("some of orders were not found")
	ErrOrderWasAlreadyIssued     = errors.New("this order was already issued")
	ErrPackageMissing            = errors.New("package is missing")
	ErrOverWeight                = "this type of package: %s can hold up to %f"
	ErrPackageDoesntExist        = "this type of package: %s does not exist"
	ErrNoRowsAffected            = errors.New("no rows affected")
	ErrNegativeLimit             = errors.New("negative limit")
)
