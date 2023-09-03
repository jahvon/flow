package parameter

import (
	"github.com/jahvon/flow/internal/backend"
)

func (p *Parameter) Data(context, masterKey string, secretBackend backend.SecretBackend) (string, error) {
	if p.Value == nil {
		return "", nil
	}

	if p.Value.Text != "" {
		return p.Value.Text, nil
	}
	ref := NormalizeKey(p.Value.Ref)
	if ref != "" {
		secret, err := secretBackend.GetSecret(context, ref)
		if err != nil {
			return "", err
		}
		if secret == "" {
			return "", nil
		}
		decryptedSecret, err := secret.Decrypt(masterKey)
		if err != nil {
			return "", err
		}
		return string(decryptedSecret), nil
	}

	return "", nil
}

func (p *Parameter) Save(context string, secretBackend backend.SecretBackend) error {
	if p.Value == nil {
		return nil
	}

	ref := NormalizeKey(p.Value.Ref)
	if ref != "" {
		secret := backend.Secret(p.Value.Text)
		err := secretBackend.SetSecret(context, ref, secret)
		if err != nil {
			return err
		}
		log.Trace().Msgf("saved secret for context %s and key %s", context, ref)
	}

	return nil
}

func (p *Parameter) Delete(context string, secretBackend backend.SecretBackend) error {
	if p.Value == nil {
		return nil
	}

	ref := NormalizeKey(p.Value.Ref)
	if ref != "" {
		err := secretBackend.DeleteSecret(context, ref)
		if err != nil {
			return err
		}
		log.Trace().Msgf("deleted secret for context %s and key %s", context, ref)
	}

	return nil
}
