package repository

import (
	"github.com/purplior/edi-adam/domain/shared/dto/pagination"
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/infra/database"
	"github.com/purplior/edi-adam/infra/database/postgre"
)

type (
	postgreRepository[T any, Q any] struct {
		client           *postgre.Client
		applyQueryOption func(
			db *postgre.DB,
			opt Q,
		) (*postgre.DB, error)
	}
)

func (r *postgreRepository[T, Q]) Read(
	session inner.Session,
	queryOption Q,
) (
	m T,
	err error,
) {
	db := r.client.DBWithContext(session)
	query, err := r.applyQueryOption(db, queryOption)
	if err != nil {
		return m, err
	}

	if err := query.First(&m).Error; err != nil {
		return m, database.ToDomainError(err)
	}

	return m, nil
}

func (r *postgreRepository[T, Q]) ReadCount(
	session inner.Session,
	queryOption Q,
) (int, error) {
	var m T
	db := r.client.DBWithContext(session)
	query, err := r.applyQueryOption(db.Model(&m), queryOption)
	if err != nil {
		return 0, err
	}

	var _count int64
	if err := query.Count(&_count).Error; err != nil {
		return 0, database.ToDomainError(err)
	}

	return int(_count), nil
}

func (r *postgreRepository[T, Q]) ReadList(
	session inner.Session,
	queryOption Q,
) (
	mArr []T,
	err error,
) {
	db := r.client.DBWithContext(session)
	query, err := r.applyQueryOption(db, queryOption)
	if err != nil {
		return nil, err
	}

	if err := query.Find(&mArr).Error; err != nil {
		return nil, database.ToDomainError(err)
	}

	return mArr, nil
}

func (r *postgreRepository[T, Q]) ReadPaginatedList(
	session inner.Session,
	option pagination.PaginationQuery[Q],
) (
	mArr []T,
	meta pagination.PaginationMeta,
	err error,
) {
	db := r.client.DBWithContext(session)
	page := option.PageRequest.Page
	pageSize := option.PageRequest.Size

	var totalCount int64
	var entity T
	countFindQuery, err := r.applyQueryOption(db.Model(&entity), option.QueryOption)
	if err != nil {
		return nil, meta, database.ToDomainError(err)
	}
	if err = countFindQuery.Count(&totalCount).Error; err != nil {
		return nil, meta, database.ToDomainError(err)
	}

	offset := (page - 1) * pageSize
	findQuery, err := r.applyQueryOption(db, option.QueryOption)
	if err != nil {
		return nil, meta, database.ToDomainError(err)
	}
	// 페이지네이션은 항상 정렬을 필수로 한다
	order := option.Order
	if len(order) == 0 {
		order = "created_at DESC"
	}

	if err = findQuery.
		Order(order).
		Offset(offset).
		Limit(pageSize).
		Find(&mArr).
		Error; err != nil {
		return nil, meta, database.ToDomainError(err)
	}

	var totalPage int = 0
	if pageSize > 0 {
		totalPage = int((totalCount + int64(pageSize) - 1) / int64(pageSize))
	}

	meta = pagination.PaginationMeta{
		Page:      page,
		Size:      pageSize,
		TotalPage: totalPage,
	}

	return mArr, meta, nil
}

func (r *postgreRepository[T, Q]) Create(
	session inner.Session,
	m T,
) (
	mRet T,
	err error,
) {
	db := r.client.DBWithContext(session)
	if err = db.Create(&m).Error; err != nil {
		return m, database.ToDomainError(err)
	}

	mRet = m

	return mRet, nil
}

func (r *postgreRepository[T, Q]) Updates(
	session inner.Session,
	queryOption Q,
	m T,
) (
	err error,
) {
	db := r.client.DBWithContext(session)
	query, err := r.applyQueryOption(db, queryOption)
	if err != nil {
		return database.ToDomainError(err)
	}

	if err = query.
		Updates(m).Error; err != nil {
		return database.ToDomainError(err)
	}

	return nil
}

func (r *postgreRepository[T, Q]) Delete(
	session inner.Session,
	queryOption Q,
) (
	err error,
) {
	db := r.client.DBWithContext(session)
	query, err := r.applyQueryOption(db, queryOption)
	if err != nil {
		return err
	}

	var entity T
	if err := query.
		Delete(&entity).Error; err != nil {
		return database.ToDomainError(err)
	}

	return nil
}
