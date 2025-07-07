package entity

type User struct {
  ID         int
  Username   string
  Name       string
  Email      *string
  Phone      string
  Mobile     string
  ImageURL   string
  Password   string
  IsActive   *bool
  RoleID     int
  CreatedAt  string
  UpdatedAt  string
}