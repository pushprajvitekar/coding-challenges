package application

import (
	"errors"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

// SignatureDeviceService implements the SignatureService interface
type SignatureDeviceService struct {
	repository     domain.SignatureDeviceRepository
	keyPairFactory domain.KeyPairFactory
	signerFactory  domain.SignerFactory
}

// NewSignatureDeviceService creates a new instance of SignatureDeviceService
func NewSignatureDeviceService(repository domain.SignatureDeviceRepository, keyPairGeneratorFactory domain.KeyPairFactory, signerFactory domain.SignerFactory) *SignatureDeviceService {
	return &SignatureDeviceService{repository: repository, keyPairFactory: keyPairGeneratorFactory, signerFactory: signerFactory}
}

// CreateSignatureDevice
func (s *SignatureDeviceService) CreateSignatureDevice(id, algorithm, label string) (*domain.SignatureDevice, error) {
	// Validate algorithm
	CryptoAlgorithm, err := domain.FromString(algorithm)
	if err != nil {
		return nil, errors.New("invalid algorithm")
	}
	// Generate key pair
	var keyPairGenerator = s.keyPairFactory.CreateKeyPairGenerator(CryptoAlgorithm)

	publicKey, privateKey, err := keyPairGenerator.GenerateKeyPair()
	if err != nil {
		return nil, err
	}

	// Create signature device
	device, err := domain.NewSignatureDevice(id, label, algorithm, publicKey, privateKey)
	if err != nil {
		return nil, err
	}

	//Persist device
	err = s.repository.AddDevice(device)
	if err != nil {
		return nil, err
	}

	return device, nil
}

func (s *SignatureDeviceService) GetSignatureDevice(id string) (*domain.SignatureDevice, error) {

	device, err := s.repository.GetDeviceByID(id)
	if err != nil {
		return nil, err
	}
	return device, nil
}
func (s *SignatureDeviceService) ListSignatureDevices() ([]*domain.SignatureDevice, error) {

	devices, err := s.repository.ListDevices()
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (s *SignatureDeviceService) SignData(id string, data string) (*domain.SignTransaction, error) {

	device, err := s.repository.GetDeviceByID(id)
	if err != nil {
		return nil, err
	}

	var signer = s.signerFactory.CreateSigner(device.KeyPairAlgorithm)

	// Use the device to sign the data
	signature, err := device.Sign(data, signer)
	if err != nil {
		return nil, err
	}
	err = s.repository.SaveSignTransaction(device, signature)
	if err != nil {
		return nil, err
	}

	return signature, nil
}
func (s *SignatureDeviceService) ListSignTransactions(deviceid string) ([]*domain.SignTransaction, error) {

	signs, err := s.repository.GetSignTransactionsForDevice(deviceid)
	if err != nil {
		return nil, err
	}
	return signs, nil
}
