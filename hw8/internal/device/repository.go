package device

type Repository interface {
    Get(num string) (Device, error)
    Create(d Device) error
    Delete(num string) error
    Update(d Device) error
}
