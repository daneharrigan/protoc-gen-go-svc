package v2

import (
	fmt "fmt"
	context "context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	is "github.com/go-ozzo/ozzo-validation/v4/is"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	publicpb "github.com/dane/protoc-gen-go-svc/example/proto/go/v2"
	privatepb "github.com/dane/protoc-gen-go-svc/example/proto/go/private"
	private "github.com/dane/protoc-gen-go-svc/example/proto/go/service/private"
)

var _ = is.Int
var _ = validation.Validate
var _ = fmt.Errorf

type Service struct {
	Validator
	Converter
	Private *private.Service
	publicpb.PeopleServer
}

func (s *Service) Create(ctx context.Context, in *publicpb.CreateRequest) (*publicpb.CreateResponse, error) {
	if err := s.ValidateCreateRequest(in); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}
	out, _, err := s.CreateImpl(ctx, in)
	return out, err
}
func (s *Service) Get(ctx context.Context, in *publicpb.GetRequest) (*publicpb.GetResponse, error) {
	if err := s.ValidateGetRequest(in); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}
	out, _, err := s.GetImpl(ctx, in)
	return out, err
}
func (s *Service) Delete(ctx context.Context, in *publicpb.DeleteRequest) (*publicpb.DeleteResponse, error) {
	if err := s.ValidateDeleteRequest(in); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}
	out, _, err := s.DeleteImpl(ctx, in)
	return out, err
}
func (s *Service) Update(ctx context.Context, in *publicpb.UpdateRequest) (*publicpb.UpdateResponse, error) {
	if err := s.ValidateUpdateRequest(in); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}
	out, _, err := s.UpdateImpl(ctx, in)
	return out, err
}
func (s *Service) CreateImpl(ctx context.Context, in *publicpb.CreateRequest, mutators ...private.CreateRequestMutator) (*publicpb.CreateResponse, *privatepb.CreateResponse, error) {
	privateIn := s.ToPrivateCreateRequest(in)
	private.ApplyCreateRequestMutators(privateIn, mutators)
	privateOut, err := s.Private.Create(ctx, privateIn)
	if err != nil {
		return nil, nil, err
	}
	out, err := s.ToPublicCreateResponse(privateOut)
	if err != nil {
		return nil, nil, status.Errorf(codes.FailedPrecondition, "%s", err)
	}
	return out, privateOut, nil
}
func (s *Service) GetImpl(ctx context.Context, in *publicpb.GetRequest, mutators ...private.FetchRequestMutator) (*publicpb.GetResponse, *privatepb.FetchResponse, error) {
	privateIn := s.ToPrivateFetchRequest(in)
	private.ApplyFetchRequestMutators(privateIn, mutators)
	privateOut, err := s.Private.Fetch(ctx, privateIn)
	if err != nil {
		return nil, nil, err
	}
	out, err := s.ToPublicGetResponse(privateOut)
	if err != nil {
		return nil, nil, status.Errorf(codes.FailedPrecondition, "%s", err)
	}
	return out, privateOut, nil
}
func (s *Service) DeleteImpl(ctx context.Context, in *publicpb.DeleteRequest, mutators ...private.DeleteRequestMutator) (*publicpb.DeleteResponse, *privatepb.DeleteResponse, error) {
	privateIn := s.ToPrivateDeleteRequest(in)
	private.ApplyDeleteRequestMutators(privateIn, mutators)
	privateOut, err := s.Private.Delete(ctx, privateIn)
	if err != nil {
		return nil, nil, err
	}
	out, err := s.ToPublicDeleteResponse(privateOut)
	if err != nil {
		return nil, nil, status.Errorf(codes.FailedPrecondition, "%s", err)
	}
	return out, privateOut, nil
}
func (s *Service) UpdateImpl(ctx context.Context, in *publicpb.UpdateRequest, mutators ...private.UpdateRequestMutator) (*publicpb.UpdateResponse, *privatepb.UpdateResponse, error) {
	privateIn := s.ToPrivateUpdateRequest(in)
	private.ApplyUpdateRequestMutators(privateIn, mutators)
	privateOut, err := s.Private.Update(ctx, privateIn)
	if err != nil {
		return nil, nil, err
	}
	out, err := s.ToPublicUpdateResponse(privateOut)
	if err != nil {
		return nil, nil, status.Errorf(codes.FailedPrecondition, "%s", err)
	}
	return out, privateOut, nil
}

const ValidatorName = "example.v2.People.Validator"

func NewValidator() Validator { return validator{} }

type Validator interface {
	Name() string
	ValidatePerson(*publicpb.Person) error
	ValidateCreateRequest(*publicpb.CreateRequest) error
	ValidateCycling(*publicpb.Cycling) error
	ValidateGetRequest(*publicpb.GetRequest) error
	ValidateDeleteRequest(*publicpb.DeleteRequest) error
	ValidateUpdateRequest(*publicpb.UpdateRequest) error
	ValidateHobby(*publicpb.Hobby) error
	ValidateHobby_Coding(*publicpb.Hobby_Coding) error
	ValidateHobby_Reading(*publicpb.Hobby_Reading) error
	ValidateHobby_Cycling(*publicpb.Hobby_Cycling) error
	ValidateCoding(*publicpb.Coding) error
	ValidateReading(*publicpb.Reading) error
}
type validator struct{}

func (v validator) Name() string { return ValidatorName }
func (v validator) ValidatePerson(in *publicpb.Person) error {
	if in == nil {
		return nil
	}
	err := validation.ValidateStruct(in,
		validation.Field(&in.Id),
		validation.Field(&in.FullName,
			validation.Required,
		),
		validation.Field(&in.Age),
		validation.Field(&in.Employment),
		validation.Field(&in.Hobby,
			validation.Required,
			validation.By(func(interface{}) error { return v.ValidateHobby(in.Hobby) }),
		),
	)
	if err != nil {
		return err
	}
	return nil
}
func (v validator) ValidateCreateRequest(in *publicpb.CreateRequest) error {
	err := validation.ValidateStruct(in,
		validation.Field(&in.Id,
			validation.Required,
			is.UUID,
		),
		validation.Field(&in.FullName,
			validation.Required,
			validation.Length(4, 0),
		),
		validation.Field(&in.Age),
		validation.Field(&in.Employment),
		validation.Field(&in.Hobby,
			validation.Required,
			validation.By(func(interface{}) error { return v.ValidateHobby(in.Hobby) }),
		),
	)
	if err != nil {
		return err
	}
	return nil
}
func (v validator) ValidateCycling(in *publicpb.Cycling) error {
	if in == nil {
		return nil
	}
	err := validation.ValidateStruct(in,
		validation.Field(&in.Style),
	)
	if err != nil {
		return err
	}
	return nil
}
func (v validator) ValidateGetRequest(in *publicpb.GetRequest) error {
	err := validation.ValidateStruct(in,
		validation.Field(&in.Id,
			validation.Required,
			is.UUID,
		),
	)
	if err != nil {
		return err
	}
	return nil
}
func (v validator) ValidateDeleteRequest(in *publicpb.DeleteRequest) error {
	err := validation.ValidateStruct(in,
		validation.Field(&in.Id),
	)
	if err != nil {
		return err
	}
	return nil
}
func (v validator) ValidateUpdateRequest(in *publicpb.UpdateRequest) error {
	err := validation.ValidateStruct(in,
		validation.Field(&in.Id,
			validation.Required,
			is.UUID,
		),
		validation.Field(&in.Person,
			validation.Required,
			validation.By(func(interface{}) error { return v.ValidatePerson(in.Person) }),
		),
	)
	if err != nil {
		return err
	}
	return nil
}
func (v validator) ValidateHobby(in *publicpb.Hobby) error {
	if in == nil {
		return nil
	}
	err := validation.ValidateStruct(in,
		validation.Field(&in.Type,
			validation.When(in.GetCoding() != nil, validation.By(func(val interface{}) error { return v.ValidateHobby_Coding(val.(*publicpb.Hobby_Coding)) })),
		),
		validation.Field(&in.Type,
			validation.When(in.GetReading() != nil, validation.By(func(val interface{}) error { return v.ValidateHobby_Reading(val.(*publicpb.Hobby_Reading)) })),
		),
		validation.Field(&in.Type,
			validation.When(in.GetCycling() != nil, validation.By(func(val interface{}) error { return v.ValidateHobby_Cycling(val.(*publicpb.Hobby_Cycling)) })),
		),
	)
	if err != nil {
		return err
	}
	return nil
}
func (v validator) ValidateHobby_Coding(in *publicpb.Hobby_Coding) error {
	if in == nil {
		return nil
	}
	err := validation.ValidateStruct(in,
		validation.Field(&in.Coding,
			validation.Required,
			validation.By(func(interface{}) error { return v.ValidateCoding(in.Coding) }),
		),
	)
	if err != nil {
		return err
	}
	return nil
}
func (v validator) ValidateHobby_Reading(in *publicpb.Hobby_Reading) error {
	if in == nil {
		return nil
	}
	err := validation.ValidateStruct(in,
		validation.Field(&in.Reading,
			validation.Required,
			validation.By(func(interface{}) error { return v.ValidateReading(in.Reading) }),
		),
	)
	if err != nil {
		return err
	}
	return nil
}
func (v validator) ValidateHobby_Cycling(in *publicpb.Hobby_Cycling) error {
	if in == nil {
		return nil
	}
	err := validation.ValidateStruct(in,
		validation.Field(&in.Cycling,
			validation.Required,
			validation.By(func(interface{}) error { return v.ValidateCycling(in.Cycling) }),
		),
	)
	if err != nil {
		return err
	}
	return nil
}
func (v validator) ValidateCoding(in *publicpb.Coding) error {
	if in == nil {
		return nil
	}
	err := validation.ValidateStruct(in,
		validation.Field(&in.Language),
	)
	if err != nil {
		return err
	}
	return nil
}
func (v validator) ValidateReading(in *publicpb.Reading) error {
	if in == nil {
		return nil
	}
	err := validation.ValidateStruct(in,
		validation.Field(&in.Genre),
	)
	if err != nil {
		return err
	}
	return nil
}

const ConverterName = "example.v2.People.Converter"

func NewConverter() Converter { return converter{} }

type Converter interface {
	Name() string
	ToPrivateCreateRequest(*publicpb.CreateRequest) *privatepb.CreateRequest
	ToPublicCreateResponse(*privatepb.CreateResponse) (*publicpb.CreateResponse, error)
	ToPrivateFetchRequest(*publicpb.GetRequest) *privatepb.FetchRequest
	ToPublicGetResponse(*privatepb.FetchResponse) (*publicpb.GetResponse, error)
	ToPrivateDeleteRequest(*publicpb.DeleteRequest) *privatepb.DeleteRequest
	ToPublicDeleteResponse(*privatepb.DeleteResponse) (*publicpb.DeleteResponse, error)
	ToPrivateUpdateRequest(*publicpb.UpdateRequest) *privatepb.UpdateRequest
	ToPublicUpdateResponse(*privatepb.UpdateResponse) (*publicpb.UpdateResponse, error)
	ToPrivatePerson(*publicpb.Person) *privatepb.Person
	ToPublicPerson(*privatepb.Person) (*publicpb.Person, error)
	ToPrivateCycling(*publicpb.Cycling) *privatepb.Cycling
	ToPublicCycling(*privatepb.Cycling) (*publicpb.Cycling, error)
	ToPrivateHobby(*publicpb.Hobby) *privatepb.Hobby
	ToPublicHobby(*privatepb.Hobby) (*publicpb.Hobby, error)
	ToPrivateCoding(*publicpb.Coding) *privatepb.Coding
	ToPublicCoding(*privatepb.Coding) (*publicpb.Coding, error)
	ToPrivateReading(*publicpb.Reading) *privatepb.Reading
	ToPublicReading(*privatepb.Reading) (*publicpb.Reading, error)
	ToPrivatePerson_Employment(publicpb.Person_Employment) privatepb.Person_Employment
	ToPublicPerson_Employment(privatepb.Person_Employment) (publicpb.Person_Employment, error)
}
type converter struct{}

func (c converter) Name() string { return ConverterName }
func (c converter) ToPrivateCreateRequest(in *publicpb.CreateRequest) *privatepb.CreateRequest {
	if in == nil {
		return nil
	}
	var out privatepb.CreateRequest
	out.Id = in.Id
	out.FullName = in.FullName
	out.Age = in.Age
	out.Employment = c.ToPrivatePerson_Employment(in.Employment)
	out.Hobby = c.ToPrivateHobby(in.Hobby)
	return &out
}
func (c converter) ToPublicCreateResponse(in *privatepb.CreateResponse) (*publicpb.CreateResponse, error) {
	if in == nil {
		return nil, nil
	}
	var required validation.Errors
	if err := required.Filter(); err != nil {
		return nil, err
	}
	var out publicpb.CreateResponse
	var err error
	out.Person, err = c.ToPublicPerson(in.Person)
	if err != nil {
		return nil, err
	}
	return &out, err
}
func (c converter) ToPrivateFetchRequest(in *publicpb.GetRequest) *privatepb.FetchRequest {
	if in == nil {
		return nil
	}
	var out privatepb.FetchRequest
	out.Id = in.Id
	return &out
}
func (c converter) ToPublicGetResponse(in *privatepb.FetchResponse) (*publicpb.GetResponse, error) {
	if in == nil {
		return nil, nil
	}
	var required validation.Errors
	if err := required.Filter(); err != nil {
		return nil, err
	}
	var out publicpb.GetResponse
	var err error
	out.Person, err = c.ToPublicPerson(in.Person)
	if err != nil {
		return nil, err
	}
	return &out, err
}
func (c converter) ToPrivateDeleteRequest(in *publicpb.DeleteRequest) *privatepb.DeleteRequest {
	if in == nil {
		return nil
	}
	var out privatepb.DeleteRequest
	out.Id = in.Id
	return &out
}
func (c converter) ToPublicDeleteResponse(in *privatepb.DeleteResponse) (*publicpb.DeleteResponse, error) {
	if in == nil {
		return nil, nil
	}
	var required validation.Errors
	if err := required.Filter(); err != nil {
		return nil, err
	}
	var out publicpb.DeleteResponse
	var err error
	return &out, err
}
func (c converter) ToPrivateUpdateRequest(in *publicpb.UpdateRequest) *privatepb.UpdateRequest {
	if in == nil {
		return nil
	}
	var out privatepb.UpdateRequest
	out.Id = in.Id
	out.Person = c.ToPrivatePerson(in.Person)
	return &out
}
func (c converter) ToPublicUpdateResponse(in *privatepb.UpdateResponse) (*publicpb.UpdateResponse, error) {
	if in == nil {
		return nil, nil
	}
	var required validation.Errors
	if err := required.Filter(); err != nil {
		return nil, err
	}
	var out publicpb.UpdateResponse
	var err error
	out.Person, err = c.ToPublicPerson(in.Person)
	if err != nil {
		return nil, err
	}
	return &out, err
}
func (c converter) ToPrivatePerson(in *publicpb.Person) *privatepb.Person {
	if in == nil {
		return nil
	}
	var out privatepb.Person
	out.Id = in.Id
	out.FullName = in.FullName
	out.Age = in.Age
	out.Employment = c.ToPrivatePerson_Employment(in.Employment)
	out.CreatedAt = in.CreatedAt
	out.UpdatedAt = in.UpdatedAt
	out.Hobby = c.ToPrivateHobby(in.Hobby)
	return &out
}
func (c converter) ToPublicPerson(in *privatepb.Person) (*publicpb.Person, error) {
	if in == nil {
		return nil, nil
	}
	var required validation.Errors
	if err := required.Filter(); err != nil {
		return nil, err
	}
	var out publicpb.Person
	var err error
	out.Id = in.Id
	out.FullName = in.FullName
	out.Age = in.Age
	out.Employment, err = c.ToPublicPerson_Employment(in.Employment)
	if err != nil {
		return nil, err
	}
	out.CreatedAt = in.CreatedAt
	out.UpdatedAt = in.UpdatedAt
	out.Hobby, err = c.ToPublicHobby(in.Hobby)
	if err != nil {
		return nil, err
	}
	return &out, err
}
func (c converter) ToPrivateCycling(in *publicpb.Cycling) *privatepb.Cycling {
	if in == nil {
		return nil
	}
	var out privatepb.Cycling
	out.Style = in.Style
	return &out
}
func (c converter) ToPublicCycling(in *privatepb.Cycling) (*publicpb.Cycling, error) {
	if in == nil {
		return nil, nil
	}
	var required validation.Errors
	if err := required.Filter(); err != nil {
		return nil, err
	}
	var out publicpb.Cycling
	var err error
	out.Style = in.Style
	return &out, err
}
func (c converter) ToPrivateHobby(in *publicpb.Hobby) *privatepb.Hobby {
	if in == nil {
		return nil
	}
	var out privatepb.Hobby
	switch in.Type.(type) {
	case *publicpb.Hobby_Coding:
		out.Type = &privatepb.Hobby_Coding{
			Coding: c.ToPrivateCoding(in.GetCoding()),
		}
	case *publicpb.Hobby_Reading:
		out.Type = &privatepb.Hobby_Reading{
			Reading: c.ToPrivateReading(in.GetReading()),
		}
	case *publicpb.Hobby_Cycling:
		out.Type = &privatepb.Hobby_Cycling{
			Cycling: c.ToPrivateCycling(in.GetCycling()),
		}
	}
	return &out
}
func (c converter) ToPublicHobby(in *privatepb.Hobby) (*publicpb.Hobby, error) {
	if in == nil {
		return nil, nil
	}
	var required validation.Errors
	if err := required.Filter(); err != nil {
		return nil, err
	}
	var out publicpb.Hobby
	var err error
	switch in.Type.(type) {
	case *privatepb.Hobby_Coding:
		var value publicpb.Hobby_Coding
		value.Coding, err = c.ToPublicCoding(in.GetCoding())
		if err == nil {
			out.Type = &value
		}
	case *privatepb.Hobby_Reading:
		var value publicpb.Hobby_Reading
		value.Reading, err = c.ToPublicReading(in.GetReading())
		if err == nil {
			out.Type = &value
		}
	case *privatepb.Hobby_Cycling:
		var value publicpb.Hobby_Cycling
		value.Cycling, err = c.ToPublicCycling(in.GetCycling())
		if err == nil {
			out.Type = &value
		}
	}
	return &out, err
}
func (c converter) ToPrivateCoding(in *publicpb.Coding) *privatepb.Coding {
	if in == nil {
		return nil
	}
	var out privatepb.Coding
	out.Language = in.Language
	return &out
}
func (c converter) ToPublicCoding(in *privatepb.Coding) (*publicpb.Coding, error) {
	if in == nil {
		return nil, nil
	}
	var required validation.Errors
	if err := required.Filter(); err != nil {
		return nil, err
	}
	var out publicpb.Coding
	var err error
	out.Language = in.Language
	return &out, err
}
func (c converter) ToPrivateReading(in *publicpb.Reading) *privatepb.Reading {
	if in == nil {
		return nil
	}
	var out privatepb.Reading
	out.Genre = in.Genre
	return &out
}
func (c converter) ToPublicReading(in *privatepb.Reading) (*publicpb.Reading, error) {
	if in == nil {
		return nil, nil
	}
	var required validation.Errors
	if err := required.Filter(); err != nil {
		return nil, err
	}
	var out publicpb.Reading
	var err error
	out.Genre = in.Genre
	return &out, err
}
func (c converter) ToPrivatePerson_Employment(in publicpb.Person_Employment) privatepb.Person_Employment {
	switch in {
	case publicpb.Person_UNSET:
		return privatepb.Person_UNDEFINED
	case publicpb.Person_FULL_TIME:
		return privatepb.Person_FULL_TIME
	case publicpb.Person_PART_TIME:
		return privatepb.Person_PART_TIME
	case publicpb.Person_UNEMPLOYED:
		return privatepb.Person_UNEMPLOYED
	}
	return privatepb.Person_UNDEFINED
}
func (c converter) ToPublicPerson_Employment(in privatepb.Person_Employment) (publicpb.Person_Employment, error) {
	switch in {
	case privatepb.Person_UNDEFINED:
		return publicpb.Person_UNSET, nil
	case privatepb.Person_FULL_TIME:
		return publicpb.Person_FULL_TIME, nil
	case privatepb.Person_PART_TIME:
		return publicpb.Person_PART_TIME, nil
	case privatepb.Person_UNEMPLOYED:
		return publicpb.Person_UNEMPLOYED, nil
	}
	return publicpb.Person_UNSET, fmt.Errorf("%q is not a supported value for this service version", in)
}
