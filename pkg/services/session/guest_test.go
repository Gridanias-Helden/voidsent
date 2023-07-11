package session

import (
	"testing"
)

func TestGuestAvatar(t *testing.T) {
	validURLs := map[string]any{
		"/assets/images/avatars/con1.png":  "",
		"/assets/images/avatars/con2.png":  "",
		"/assets/images/avatars/con3.png":  "",
		"/assets/images/avatars/con4.png":  "",
		"/assets/images/avatars/con5.png":  "",
		"/assets/images/avatars/con6.png":  "",
		"/assets/images/avatars/con7.png":  "",
		"/assets/images/avatars/con8.png":  "",
		"/assets/images/avatars/con9.png":  "",
		"/assets/images/avatars/con10.png": "",
		"/assets/images/avatars/con11.png": "",
		"/assets/images/avatars/con12.png": "",
		"/assets/images/avatars/con13.png": "",
		"/assets/images/avatars/con14.png": "",
		"/assets/images/avatars/con15.png": "",
		"/assets/images/avatars/con16.png": "",
		"/assets/images/avatars/con17.png": "",
		"/assets/images/avatars/con18.png": "",
		"/assets/images/avatars/con19.png": "",
		"/assets/images/avatars/con20.png": "",
		"/assets/images/avatars/con21.png": "",
		"/assets/images/avatars/con22.png": "",
		"/assets/images/avatars/con23.png": "",
		"/assets/images/avatars/con24.png": "",
		"/assets/images/avatars/con25.png": "",
		"/assets/images/avatars/con26.png": "",
		"/assets/images/avatars/con27.png": "",
		"/assets/images/avatars/con28.png": "",
		"/assets/images/avatars/con29.png": "",
		"/assets/images/avatars/con30.png": "",
		"/assets/images/avatars/con31.png": "",
		"/assets/images/avatars/con32.png": "",
		"/assets/images/avatars/con33.png": "",
		"/assets/images/avatars/con34.png": "",
		"/assets/images/avatars/con35.png": "",
		"/assets/images/avatars/con36.png": "",
		"/assets/images/avatars/con37.png": "",
		"/assets/images/avatars/con38.png": "",
		"/assets/images/avatars/con39.png": "",
		"/assets/images/avatars/con40.png": "",
		"/assets/images/avatars/con41.png": "",
		"/assets/images/avatars/con42.png": "",
	}

	for index := 0; index < 1000000; index++ {
		uri := GuestAvatar()
		if _, ok := validURLs[uri]; !ok {
			t.Errorf("URI %q is not valid", uri)
		}
	}
}
