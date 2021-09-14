package service

import (
	"testing"
)

func TestValidateURL(t *testing.T) {
	url := "https://www.istockphoto.com/en/search/2/image?mediatype=photography&phrase=marguerite&utm_source=pixabay&utm_medium=affiliate&utm_campaign=ADP_photo_sponsored&referrer_url=http:%2F%2Fpixabay.com%2Fphotos%2Fmarguerite-daisy-flower-white-729510%2F&utm_term=marguerite"
	expected := bool
	isValid, err := validate(url)

	if !isValid || err != nil {
		t.Errorf("validate(url) = %v; want true", isValid)
	}
}
