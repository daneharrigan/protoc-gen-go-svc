package v2

import (
	"context"
	publicpb "github.com/dane/protoc-gen-go-svc/example/proto/go/v2"
	privatepb "github.com/dane/protoc-gen-go-svc/example/proto/go/private"
	private "github.com/dane/protoc-gen-go-svc/example/proto/go/service/private"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	is "github.com/go-ozzo/ozzo-validation/v4/is"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

var _ = is.Int

const ValidatorName = "example.v2.People.Validator"

type Validator interface {
	Name() string
	ValidateCreateRequest(*publicpb.CreateRequest) error
	ValidateGetRequest(*publicpb.GetRequest) error
	ValidateDeleteRequest(*publicpb.DeleteRequest) error
}

func NewValidator() Validator { return validator{} }

type validator struct{}

func (v validator) Name() string { return ValidatorName }
func (v validator) ValidateCreateRequest(in *publicpb.CreateRequest) error {
	err := validation.ValidateStruct(in)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
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
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return nil
}
func (v validator) ValidateDeleteRequest(in *publicpb.DeleteRequest) error {
	err := validation.ValidateStruct(in)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
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
	ToPrivatePerson(*publicpb.Person) *privatepb.Person
	ToPublicPerson(*privatepb.Person) (*publicpb.Person, error)
	ToPrivatePerson_Employment(publicpb.Person_Employment) privatepb.Person_Employment
	ToPublicPerson_Employment(privatepb.Person_Employment) (publicpb.Person_Employment, error)
}
type converter struct{}

func (c converter) Name() string { return ConverterName }
func (c converter) ToPrivateCreateRequest(in *publicpb.CreateRequest) *privatepb.CreateRequest {
	var out privatepb.CreateRequest
	out.FullName = in.FullName
	out.Age = in.Age
	out.Employment = c.ToPrivatePerson_Employment(in.Employment)
	return &out
}
func (c converter) ToPublicCreateResponse(in *privatepb.CreateResponse) (*publicpb.CreateResponse, error) {
	var out publicpb.CreateResponse
	var err error
	out.Person, err = c.ToPublicPerson(in.Person)
	if err != nil {
		return nil, err
	}
	return &out, err
}
func (c converter) ToPrivateFetchRequest(in *publicpb.GetRequest) *privatepb.FetchRequest {
	var out privatepb.FetchRequest
	out.Id = in.Id
	return &out
}
func (c converter) ToPublicGetResponse(in *privatepb.FetchResponse) (*publicpb.GetResponse, error) {
	var out publicpb.GetResponse
	var err error
	out.Person, err = c.ToPublicPerson(in.Person)
	if err != nil {
		return nil, err
	}
	return &out, err
}
func (c converter) ToPrivateDeleteRequest(in *publicpb.DeleteRequest) *privatepb.DeleteRequest {
	var out privatepb.DeleteRequest
	out.Id = in.Id
	return &out
}
func (c converter) ToPublicDeleteResponse(in *privatepb.DeleteResponse) (*publicpb.DeleteResponse, error) {
	var out publicpb.DeleteResponse
	var err error
	return &out, err
}
func (c converter) ToPrivatePerson(in *publicpb.Person) *privatepb.Person {
	var out privatepb.Person
	out.Id = in.Id
	out.FullName = in.FullName
	out.Age = in.Age
	out.Employment = c.ToPrivatePerson_Employment(in.Employment)
	out.CreatedAt = in.CreatedAt
	out.UpdatedAt = in.UpdatedAt
	return &out
}
func (c converter) ToPublicPerson(in *privatepb.Person) (*publicpb.Person, error) {
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
	return publicpb.Person_UNSET, status.Errorf(codes.FailedPrecondition, "unexpected value %q", in)
}

type Service struct {
	publicpb.PeopleServer
	Validator
	Converter
	Private *private.Service
}

func (s *Service) Create(ctx context.Context, in *publicpb.CreateRequest) (*publicpb.CreateResponse, error) {
	if err := s.ValidateCreateRequest(in); err != nil {
		return nil, err
	}
	out, _, err := s.CreateImpl(ctx, in)
	return out, err
}
func (s *Service) CreateImpl(ctx context.Context, in *publicpb.CreateRequest, mutators ...private.CreateRequestMutator) (*publicpb.CreateResponse, *privatepb.CreateResponse, error) {
	privIn := s.ToPrivateCreateRequest(in)
	privOut, err := s.Private.Create(ctx, privIn)
	if err != nil {
		return nil, nil, err
	}
	out, err := s.ToPublicCreateResponse(privOut)
	if err != nil {
		return nil, nil, err
	}
	return out, privOut, err
}
func (s *Service) Get(ctx context.Context, in *publicpb.GetRequest) (*publicpb.GetResponse, error) {
	if err := s.ValidateGetRequest(in); err != nil {
		return nil, err
	}
	out, _, err := s.GetImpl(ctx, in)
	return out, err
}
func (s *Service) GetImpl(ctx context.Context, in *publicpb.GetRequest, mutators ...private.FetchRequestMutator) (*publicpb.GetResponse, *privatepb.FetchResponse, error) {
	privIn := s.ToPrivateFetchRequest(in)
	privOut, err := s.Private.Fetch(ctx, privIn)
	if err != nil {
		return nil, nil, err
	}
	out, err := s.ToPublicGetResponse(privOut)
	if err != nil {
		return nil, nil, err
	}
	return out, privOut, err
}
func (s *Service) Delete(ctx context.Context, in *publicpb.DeleteRequest) (*publicpb.DeleteResponse, error) {
	if err := s.ValidateDeleteRequest(in); err != nil {
		return nil, err
	}
	out, _, err := s.DeleteImpl(ctx, in)
	return out, err
}
func (s *Service) DeleteImpl(ctx context.Context, in *publicpb.DeleteRequest, mutators ...private.DeleteRequestMutator) (*publicpb.DeleteResponse, *privatepb.DeleteResponse, error) {
	privIn := s.ToPrivateDeleteRequest(in)
	privOut, err := s.Private.Delete(ctx, privIn)
	if err != nil {
		return nil, nil, err
	}
	out, err := s.ToPublicDeleteResponse(privOut)
	if err != nil {
		return nil, nil, err
	}
	return out, privOut, err
}
