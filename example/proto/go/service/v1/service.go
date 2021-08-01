package v1

import (
	fmt "fmt"
	context "context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	is "github.com/go-ozzo/ozzo-validation/v4/is"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	publicpb "github.com/dane/protoc-gen-go-svc/example/proto/go/v1"
	nextpb "github.com/dane/protoc-gen-go-svc/example/proto/go/v2"
	privatepb "github.com/dane/protoc-gen-go-svc/example/proto/go/private"
	private "github.com/dane/protoc-gen-go-svc/example/proto/go/service/private"
	next "github.com/dane/protoc-gen-go-svc/example/proto/go/service/v2"
)

var _ = is.Int
var _ = validation.Validate
var _ = fmt.Errorf

type Service struct {
	Validator
	Converter
	Private *private.Service
	Next    *next.Service
	publicpb.PeopleServer
}

const ValidatorName = "example.v1.People.Validator"

func NewValidator() Validator { return validator{} }

type Validator interface {
	Name() string
	ValidateBiking(*publicpb.Biking) error
	ValidateCoding(*publicpb.Coding) error
	ValidateCreateRequest(*publicpb.CreateRequest) error
	ValidateDeleteRequest(*publicpb.DeleteRequest) error
	ValidateGetRequest(*publicpb.GetRequest) error
	ValidateHobby(*publicpb.Hobby) error
	ValidateHobby_Coding(*publicpb.Hobby_Coding) error
	ValidateHobby_Reading(*publicpb.Hobby_Reading) error
	ValidateHobby_Biking(*publicpb.Hobby_Biking) error
	ValidateListRequest(*publicpb.ListRequest) error
	ValidatePerson(*publicpb.Person) error
	ValidateReading(*publicpb.Reading) error
}
type validator struct{}

func (v validator) Name() string { return ValidatorName }
func (v validator) ValidateBiking(in *publicpb.Biking) error {
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
func (v validator) ValidateCreateRequest(in *publicpb.CreateRequest) error {
	err := validation.ValidateStruct(in,
		validation.Field(&in.Id,
			validation.Required,
			is.UUID,
		),
		validation.Field(&in.FirstName,
			validation.Required,
			validation.Length(2, 0),
		),
		validation.Field(&in.LastName,
			validation.Required,
			validation.Length(2, 0),
		),
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
func (v validator) ValidateDeleteRequest(in *publicpb.DeleteRequest) error {
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
			validation.When(in.GetBiking() != nil, validation.By(func(val interface{}) error { return v.ValidateHobby_Biking(val.(*publicpb.Hobby_Biking)) })),
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
func (v validator) ValidateHobby_Biking(in *publicpb.Hobby_Biking) error {
	if in == nil {
		return nil
	}
	err := validation.ValidateStruct(in,
		validation.Field(&in.Biking,
			validation.Required,
			validation.By(func(interface{}) error { return v.ValidateBiking(in.Biking) }),
		),
	)
	if err != nil {
		return err
	}
	return nil
}
func (v validator) ValidateListRequest(in *publicpb.ListRequest) error {
	err := validation.ValidateStruct(in)
	if err != nil {
		return err
	}
	return nil
}
func (v validator) ValidatePerson(in *publicpb.Person) error {
	if in == nil {
		return nil
	}
	err := validation.ValidateStruct(in,
		validation.Field(&in.Id),
		validation.Field(&in.FirstName),
		validation.Field(&in.LastName),
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

const ConverterName = "example.v1.People.Converter"

func NewConverter() Converter { return converter{} }

type Converter interface {
	Name() string
	ToNextCreateRequest(*publicpb.CreateRequest) *nextpb.CreateRequest
	ToPublicCreateResponse(*nextpb.CreateResponse, *privatepb.CreateResponse) (*publicpb.CreateResponse, error)
	ToNextDeleteRequest(*publicpb.DeleteRequest) *nextpb.DeleteRequest
	ToPublicDeleteResponse(*nextpb.DeleteResponse, *privatepb.DeleteResponse) (*publicpb.DeleteResponse, error)
	ToNextGetRequest(*publicpb.GetRequest) *nextpb.GetRequest
	ToPublicGetResponse(*nextpb.GetResponse, *privatepb.FetchResponse) (*publicpb.GetResponse, error)
	ToPrivateListRequest(*publicpb.ListRequest) *privatepb.ListRequest
	ToNextCycling(*publicpb.Biking) *nextpb.Cycling
	ToPublicBiking(*nextpb.Cycling, *privatepb.Cycling) (*publicpb.Biking, error)
	ToNextCoding(*publicpb.Coding) *nextpb.Coding
	ToPublicCoding(*nextpb.Coding, *privatepb.Coding) (*publicpb.Coding, error)
	ToNextHobby(*publicpb.Hobby) *nextpb.Hobby
	ToPublicHobby(*nextpb.Hobby, *privatepb.Hobby) (*publicpb.Hobby, error)
	ToNextPerson(*publicpb.Person) *nextpb.Person
	ToPublicPerson(*nextpb.Person, *privatepb.Person) (*publicpb.Person, error)
	ToNextReading(*publicpb.Reading) *nextpb.Reading
	ToPublicReading(*nextpb.Reading, *privatepb.Reading) (*publicpb.Reading, error)
	ToNextPerson_Employment(publicpb.Person_Employment) nextpb.Person_Employment
	ToPublicPerson_Employment(nextpb.Person_Employment) (publicpb.Person_Employment, error)
	ToDeprecatedPublicListResponse(*privatepb.ListResponse) (*publicpb.ListResponse, error)
	ToDeprecatedPublicBiking(*privatepb.Cycling) (*publicpb.Biking, error)
	ToDeprecatedPublicCoding(*privatepb.Coding) (*publicpb.Coding, error)
	ToDeprecatedPublicHobby(*privatepb.Hobby) (*publicpb.Hobby, error)
	ToDeprecatedPublicPerson(*privatepb.Person) (*publicpb.Person, error)
	ToDeprecatedPublicReading(*privatepb.Reading) (*publicpb.Reading, error)
	ToDeprecatedPublicPerson_Employment(privatepb.Person_Employment) (publicpb.Person_Employment, error)
}
type converter struct{}

func (c converter) Name() string { return ConverterName }
func (c converter) ToNextCreateRequest(in *publicpb.CreateRequest) *nextpb.CreateRequest {
	if in == nil {
		return nil
	}
	var out nextpb.CreateRequest
	out.Id = in.Id
	out.Employment = c.ToNextPerson_Employment(in.Employment)
	out.Hobby = c.ToNextHobby(in.Hobby)
	return &out
}
func (c converter) ToPublicCreateResponse(nextIn *nextpb.CreateResponse, privateIn *privatepb.CreateResponse) (*publicpb.CreateResponse, error) {
	if nextIn == nil || privateIn == nil {
		return nil, nil
	}
	required := validation.Errors{}
	if err := required.Filter(); err != nil {
		return nil, err
	}
	var out publicpb.CreateResponse
	var err error
	out.Person, err = c.ToPublicPerson(nextIn.Person, privateIn.Person)
	if err != nil {
		return nil, err
	}
	return &out, err
}
func (c converter) ToNextDeleteRequest(in *publicpb.DeleteRequest) *nextpb.DeleteRequest {
	if in == nil {
		return nil
	}
	var out nextpb.DeleteRequest
	out.Id = in.Id
	return &out
}
func (c converter) ToPublicDeleteResponse(nextIn *nextpb.DeleteResponse, privateIn *privatepb.DeleteResponse) (*publicpb.DeleteResponse, error) {
	if nextIn == nil || privateIn == nil {
		return nil, nil
	}
	required := validation.Errors{}
	if err := required.Filter(); err != nil {
		return nil, err
	}
	var out publicpb.DeleteResponse
	var err error
	return &out, err
}
func (c converter) ToNextGetRequest(in *publicpb.GetRequest) *nextpb.GetRequest {
	if in == nil {
		return nil
	}
	var out nextpb.GetRequest
	out.Id = in.Id
	return &out
}
func (c converter) ToPublicGetResponse(nextIn *nextpb.GetResponse, privateIn *privatepb.FetchResponse) (*publicpb.GetResponse, error) {
	if nextIn == nil || privateIn == nil {
		return nil, nil
	}
	required := validation.Errors{}
	if err := required.Filter(); err != nil {
		return nil, err
	}
	var out publicpb.GetResponse
	var err error
	out.Person, err = c.ToPublicPerson(nextIn.Person, privateIn.Person)
	if err != nil {
		return nil, err
	}
	return &out, err
}
func (c converter) ToPrivateListRequest(in *publicpb.ListRequest) *privatepb.ListRequest {
	if in == nil {
		return nil
	}
	var out privatepb.ListRequest
	return &out
}
func (c converter) ToNextCycling(in *publicpb.Biking) *nextpb.Cycling {
	if in == nil {
		return nil
	}
	var out nextpb.Cycling
	out.Style = in.Style
	return &out
}
func (c converter) ToPublicBiking(nextIn *nextpb.Cycling, privateIn *privatepb.Cycling) (*publicpb.Biking, error) {
	if nextIn == nil || privateIn == nil {
		return nil, nil
	}
	required := validation.Errors{}
	if err := required.Filter(); err != nil {
		return nil, err
	}
	var out publicpb.Biking
	var err error
	out.Style = nextIn.Style
	return &out, err
}
func (c converter) ToNextCoding(in *publicpb.Coding) *nextpb.Coding {
	if in == nil {
		return nil
	}
	var out nextpb.Coding
	out.Language = in.Language
	return &out
}
func (c converter) ToPublicCoding(nextIn *nextpb.Coding, privateIn *privatepb.Coding) (*publicpb.Coding, error) {
	if nextIn == nil || privateIn == nil {
		return nil, nil
	}
	required := validation.Errors{}
	if err := required.Filter(); err != nil {
		return nil, err
	}
	var out publicpb.Coding
	var err error
	out.Language = nextIn.Language
	return &out, err
}
func (c converter) ToNextHobby(in *publicpb.Hobby) *nextpb.Hobby {
	if in == nil {
		return nil
	}
	var out nextpb.Hobby
	switch in.Type.(type) {
	case *publicpb.Hobby_Coding:
		out.Type = &nextpb.Hobby_Coding{
			Coding: c.ToNextCoding(in.GetCoding()),
		}
	case *publicpb.Hobby_Reading:
		out.Type = &nextpb.Hobby_Reading{
			Reading: c.ToNextReading(in.GetReading()),
		}
	case *publicpb.Hobby_Biking:
		out.Type = &nextpb.Hobby_Cycling{
			Cycling: c.ToNextCycling(in.GetBiking()),
		}
	}
	return &out
}
func (c converter) ToPublicHobby(nextIn *nextpb.Hobby, privateIn *privatepb.Hobby) (*publicpb.Hobby, error) {
	if nextIn == nil || privateIn == nil {
		return nil, nil
	}
	required := validation.Errors{}
	if err := required.Filter(); err != nil {
		return nil, err
	}
	var out publicpb.Hobby
	var err error
	switch nextIn.Type.(type) {
	case *nextpb.Hobby_Coding:
		var value publicpb.Hobby_Coding
		value.Coding, err = c.ToPublicCoding(nextIn.GetCoding(), privateIn.GetCoding())
		out.Type = &value
	case *nextpb.Hobby_Reading:
		var value publicpb.Hobby_Reading
		value.Reading, err = c.ToPublicReading(nextIn.GetReading(), privateIn.GetReading())
		out.Type = &value
	case *nextpb.Hobby_Cycling:
		var value publicpb.Hobby_Biking
		value.Biking, err = c.ToPublicBiking(nextIn.GetCycling(), privateIn.GetCycling())
		out.Type = &value
	}
	return &out, err
}
func (c converter) ToNextPerson(in *publicpb.Person) *nextpb.Person {
	if in == nil {
		return nil
	}
	var out nextpb.Person
	out.Id = in.Id
	out.Employment = c.ToNextPerson_Employment(in.Employment)
	out.CreatedAt = in.CreatedAt
	out.UpdatedAt = in.UpdatedAt
	out.Hobby = c.ToNextHobby(in.Hobby)
	return &out
}
func (c converter) ToPublicPerson(nextIn *nextpb.Person, privateIn *privatepb.Person) (*publicpb.Person, error) {
	if nextIn == nil || privateIn == nil {
		return nil, nil
	}
	required := validation.Errors{}
	required["Id"] = validation.Validate(nextIn.Id, validation.Required)
	required["FirstName"] = validation.Validate(privateIn.FirstName, validation.Required)
	required["LastName"] = validation.Validate(privateIn.LastName, validation.Required)
	required["Employment"] = validation.Validate(nextIn.Employment, validation.Required)
	if err := required.Filter(); err != nil {
		return nil, err
	}
	var out publicpb.Person
	var err error
	out.Id = nextIn.Id
	out.FirstName = privateIn.FirstName
	out.LastName = privateIn.LastName
	out.Employment, err = c.ToPublicPerson_Employment(nextIn.Employment)
	if err != nil {
		return nil, err
	}
	out.CreatedAt = nextIn.CreatedAt
	out.UpdatedAt = nextIn.UpdatedAt
	out.Hobby, err = c.ToPublicHobby(nextIn.Hobby, privateIn.Hobby)
	if err != nil {
		return nil, err
	}
	return &out, err
}
func (c converter) ToNextReading(in *publicpb.Reading) *nextpb.Reading {
	if in == nil {
		return nil
	}
	var out nextpb.Reading
	out.Genre = in.Genre
	return &out
}
func (c converter) ToPublicReading(nextIn *nextpb.Reading, privateIn *privatepb.Reading) (*publicpb.Reading, error) {
	if nextIn == nil || privateIn == nil {
		return nil, nil
	}
	required := validation.Errors{}
	if err := required.Filter(); err != nil {
		return nil, err
	}
	var out publicpb.Reading
	var err error
	out.Genre = nextIn.Genre
	return &out, err
}
func (c converter) ToNextPerson_Employment(in publicpb.Person_Employment) nextpb.Person_Employment {
	switch in {
	case publicpb.Person_UNSET:
		return nextpb.Person_UNSET
	case publicpb.Person_EMPLOYED:
		return nextpb.Person_FULL_TIME
	case publicpb.Person_UNEMPLOYED:
		return nextpb.Person_UNEMPLOYED
	}
	return nextpb.Person_UNSET
}
func (c converter) ToPublicPerson_Employment(in nextpb.Person_Employment) (publicpb.Person_Employment, error) {
	switch in {
	case nextpb.Person_UNSET:
		return publicpb.Person_UNSET, nil
	case nextpb.Person_FULL_TIME:
		return publicpb.Person_EMPLOYED, nil
	case nextpb.Person_PART_TIME:
		return publicpb.Person_EMPLOYED, nil
	case nextpb.Person_UNEMPLOYED:
		return publicpb.Person_UNEMPLOYED, nil
	}
	return publicpb.Person_UNSET, fmt.Errorf("%q is not a supported value for this service version", in)
}
func (c converter) ToDeprecatedPublicListResponse(in *privatepb.ListResponse) (*publicpb.ListResponse, error) {
	if in == nil {
		return nil, nil
	}
	var required validation.Errors
	if err := required.Filter(); err != nil {
		return nil, err
	}
	var out publicpb.ListResponse
	var err error
	for _, item := range in.People {
		conv, err := c.ToDeprecatedPublicPerson(item)
		if err != nil {
			return nil, err
		}
		out.People = append(out.People, conv)
	}
	return &out, err
}
func (c converter) ToDeprecatedPublicBiking(in *privatepb.Cycling) (*publicpb.Biking, error) {
	if in == nil {
		return nil, nil
	}
	var required validation.Errors
	if err := required.Filter(); err != nil {
		return nil, err
	}
	var out publicpb.Biking
	var err error
	out.Style = in.Style
	return &out, err
}
func (c converter) ToDeprecatedPublicCoding(in *privatepb.Coding) (*publicpb.Coding, error) {
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
func (c converter) ToDeprecatedPublicHobby(in *privatepb.Hobby) (*publicpb.Hobby, error) {
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
		value.Coding, err = c.ToDeprecatedPublicCoding(in.GetCoding())
		if err == nil {
			out.Type = &value
		}
	case *privatepb.Hobby_Reading:
		var value publicpb.Hobby_Reading
		value.Reading, err = c.ToDeprecatedPublicReading(in.GetReading())
		if err == nil {
			out.Type = &value
		}
	case *privatepb.Hobby_Cycling:
		var value publicpb.Hobby_Biking
		value.Biking, err = c.ToDeprecatedPublicBiking(in.GetCycling())
		if err == nil {
			out.Type = &value
		}
	}
	return &out, err
}
func (c converter) ToDeprecatedPublicPerson(in *privatepb.Person) (*publicpb.Person, error) {
	if in == nil {
		return nil, nil
	}
	var required validation.Errors
	required["Id"] = validation.Validate(in.Id, validation.Required)
	required["FirstName"] = validation.Validate(in.FirstName, validation.Required)
	required["LastName"] = validation.Validate(in.LastName, validation.Required)
	required["Employment"] = validation.Validate(in.Employment, validation.Required)
	if err := required.Filter(); err != nil {
		return nil, err
	}
	var out publicpb.Person
	var err error
	out.Id = in.Id
	out.FirstName = in.FirstName
	out.LastName = in.LastName
	out.Employment, err = c.ToDeprecatedPublicPerson_Employment(in.Employment)
	if err != nil {
		return nil, err
	}
	out.CreatedAt = in.CreatedAt
	out.UpdatedAt = in.UpdatedAt
	out.Hobby, err = c.ToDeprecatedPublicHobby(in.Hobby)
	if err != nil {
		return nil, err
	}
	return &out, err
}
func (c converter) ToDeprecatedPublicReading(in *privatepb.Reading) (*publicpb.Reading, error) {
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
func (c converter) ToDeprecatedPublicPerson_Employment(in privatepb.Person_Employment) (publicpb.Person_Employment, error) {
	switch in {
	case privatepb.Person_UNDEFINED:
		return publicpb.Person_UNSET, nil
	case privatepb.Person_FULL_TIME:
		return publicpb.Person_EMPLOYED, nil
	case privatepb.Person_PART_TIME:
		return publicpb.Person_EMPLOYED, nil
	case privatepb.Person_UNEMPLOYED:
		return publicpb.Person_UNEMPLOYED, nil
	}
	return publicpb.Person_UNSET, fmt.Errorf("%q is not a supported value for this service version", in)
}
func (s *Service) Create(ctx context.Context, in *publicpb.CreateRequest) (*publicpb.CreateResponse, error) {
	if err := s.ValidateCreateRequest(in); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}
	out, _, err := s.CreateImpl(ctx, in)
	return out, err
}
func (s *Service) Delete(ctx context.Context, in *publicpb.DeleteRequest) (*publicpb.DeleteResponse, error) {
	if err := s.ValidateDeleteRequest(in); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}
	out, _, err := s.DeleteImpl(ctx, in)
	return out, err
}
func (s *Service) Get(ctx context.Context, in *publicpb.GetRequest) (*publicpb.GetResponse, error) {
	if err := s.ValidateGetRequest(in); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}
	out, _, err := s.GetImpl(ctx, in)
	return out, err
}
func (s *Service) List(ctx context.Context, in *publicpb.ListRequest) (*publicpb.ListResponse, error) {
	if err := s.ValidateListRequest(in); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}
	out, _, err := s.ListImpl(ctx, in)
	return out, err
}
func (s *Service) CreateImpl(ctx context.Context, in *publicpb.CreateRequest, mutators ...private.CreateRequestMutator) (*publicpb.CreateResponse, *privatepb.CreateResponse, error) {
	mutators = append(mutators, private.SetCreateRequest_FirstName(in.FirstName))
	mutators = append(mutators, private.SetCreateRequest_LastName(in.LastName))
	nextIn := s.ToNextCreateRequest(in)
	nextOut, privateOut, err := s.Next.CreateImpl(ctx, nextIn, mutators...)
	if err != nil {
		return nil, nil, err
	}
	out, err := s.ToPublicCreateResponse(nextOut, privateOut)
	if err != nil {
		return nil, nil, status.Errorf(codes.FailedPrecondition, "%s", err)
	}
	return out, privateOut, nil
}
func (s *Service) DeleteImpl(ctx context.Context, in *publicpb.DeleteRequest, mutators ...private.DeleteRequestMutator) (*publicpb.DeleteResponse, *privatepb.DeleteResponse, error) {
	nextIn := s.ToNextDeleteRequest(in)
	nextOut, privateOut, err := s.Next.DeleteImpl(ctx, nextIn, mutators...)
	if err != nil {
		return nil, nil, err
	}
	out, err := s.ToPublicDeleteResponse(nextOut, privateOut)
	if err != nil {
		return nil, nil, status.Errorf(codes.FailedPrecondition, "%s", err)
	}
	return out, privateOut, nil
}
func (s *Service) GetImpl(ctx context.Context, in *publicpb.GetRequest, mutators ...private.FetchRequestMutator) (*publicpb.GetResponse, *privatepb.FetchResponse, error) {
	nextIn := s.ToNextGetRequest(in)
	nextOut, privateOut, err := s.Next.GetImpl(ctx, nextIn, mutators...)
	if err != nil {
		return nil, nil, err
	}
	out, err := s.ToPublicGetResponse(nextOut, privateOut)
	if err != nil {
		return nil, nil, status.Errorf(codes.FailedPrecondition, "%s", err)
	}
	return out, privateOut, nil
}
func (s *Service) ListImpl(ctx context.Context, in *publicpb.ListRequest, mutators ...private.ListRequestMutator) (*publicpb.ListResponse, *privatepb.ListResponse, error) {
	privateIn := s.ToPrivateListRequest(in)
	private.ApplyListRequestMutators(privateIn, mutators)
	privateOut, err := s.Private.List(ctx, privateIn)
	if err != nil {
		return nil, nil, err
	}
	out, err := s.ToDeprecatedPublicListResponse(privateOut)
	if err != nil {
		return nil, nil, status.Errorf(codes.FailedPrecondition, "%s", err)
	}
	return out, privateOut, nil
}
