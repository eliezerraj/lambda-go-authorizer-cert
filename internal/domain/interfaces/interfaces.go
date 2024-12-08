package interfaces

type Cache interface {
	VerifyCertCRL() (error)
}