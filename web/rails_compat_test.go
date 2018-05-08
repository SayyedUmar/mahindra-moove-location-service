package web

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestDecodeRailsSessionShouldUnescapeUrlEncoding(t *testing.T) {
	keyBase := "daef0f128e8c899673cbcc71b3d4854112006de625db4bab0116818c7a7dc5cae52fc531dfb5f153d1b7fb3021331de744683b0346317751115c1837cbc69672"
	cookie := "OGFNYmtFejZTZDc0YkpXTjFXaUp1NHZjaHRpcjVKajM3eFZPWEZjamlaTk5UT0huY1R5YUlFdndITHUydTlzVUNDNEZJY0FGeUN5Z0s2Um52ODJqekQ4bWxYaHI5SmoreVBHYUtkenJpSm1KdWRrMXYrZmc2a3FXTlNZYUlSRy8vSTc4bStwM2tHVDFGZ1NrV3V6cVIxa2dnWUNMWE5FNXFMamNjVzBpZHhpajBXaGJlZEN4QmV0bEJkUENSZG0vSmJ3RzVMdUNGakRyQUU3blh5d01uclFVM1Y5Z244SGxEWkhPZk95TUdyaz0tLWVJZjR6bU5Uck5lTCtzT3V1VXRJdWc9PQ%3D%3D--40a7c5ea0818daabde304fd65a878aaef534cd27"
	var sess map[string]interface{}
	decrypted := DecodeRailsSession(cookie, keyBase)
	buffer := bytes.NewBuffer([]byte(decrypted))
	decoder := json.NewDecoder(buffer)
	err := decoder.Decode(&sess)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if sess["session_id"].(string) != "dfa336b7bc2bb7e84bef023fd0291aea" {
		t.Logf("expected %s found %s\n", "dfa336b7bc2bb7e84bef023fd0291aea", sess["session_id"].(string))
		t.FailNow()
	}
}

func TestDecodeRailsSessionShouldDecodeSession(t *testing.T) {

	keyBase := "f1d186616befd0912ed643cdc621377baa17368402970cfca9eaaf75f93286da121c22f1576ac5399a0d4c9ab3026849ebb67cd617437d73835c136e1c40a946"
	// this is the decrypted session content
	// {"session_id":"d8b8304f5f339a818e127aca2dfab742","user_return_to":"/","warden.user.user.key":[[4],"$2a$11$gZkXn2gGS11ROQQJ1./hGO"],"_csrf_token":"atHkScP0CWcPrcxIJtdPk2Yg1aKhPTQ5HYg+sP/rjts="}
	cookie := "bVdpR2NLeUhBTXFMUk5rdUMrUWtGTnlrWDhoNCtQZmFMRVVPUTJNQlhkR2VXOU9oME1vZ1NabHBqREFNbVAzYkVnWCtSUjRCaXJaZjBIbEFodXl6Y28yV0IxSmQ1bHgzOTJoNlZQQzN2TzZsSnNYbUgzWkFMb291Q3FRTWozVmc3elNxSi9LTUN6STA3dnk4bnRFZDRUUU94K2VteUIwNkUxdWF0Zk8wb2x3a3h4OWw1Q3BhYWhGTGZDSFJDdjdUL3lwRi9URVNMUEhVOGtSN3dPUHJuTkdKTzdyTnMzcDlaUHVxNzdQVTh4aHo1ZFVUWkJwdWY4M0tKZVE2THpMdzFiM1FEYU13dlh3dTZGOFMyWDF2UXc9PS0tU1l5ZjhrVnN3NTYyaWxnZkZMZFNIdz09"

	var sess map[string]interface{}

	decrypted := DecodeRailsSession(cookie, keyBase)
	buffer := bytes.NewBuffer([]byte(decrypted))
	decoder := json.NewDecoder(buffer)
	err := decoder.Decode(&sess)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	if sess["session_id"].(string) != "d8b8304f5f339a818e127aca2dfab742" {
		t.Logf("expected %s found %s\n", "d8b8304f5f339a818e127aca2dfab742", sess["session_id"].(string))
		t.FailNow()
	}
}

func TestExtractUserIdShouldExtractUserId(t *testing.T) {
	session := `{"session_id":"d8b8304f5f339a818e127aca2dfab742","user_return_to":"/","warden.user.user.key":[[42],"$2a$11$gZkXn2gGS11ROQQJ1./hGO"],"_csrf_token":"atHkScP0CWcPrcxIJtdPk2Yg1aKhPTQ5HYg+sP/rjts="}`
	uid, err := ExtractUserId(session)
	if err != nil {
		t.Log("expected no errors found error")
		t.Log(err)
		t.FailNow()
	}
	if uid != 42 {
		t.Logf("expected userID to be 42 found %d\n", uid)
		t.FailNow()
	}
}
func TestExtractUserIdShouldFailGraceFully(t *testing.T) {
	session := `{"session_id":"d8b8304f5f339a818e127aca2dfab742","user_return_to":"/","warden.user.user.key":[42,"$2a$11$gZkXn2gGS11ROQQJ1./hGO"],"_csrf_token":"atHkScP0CWcPrcxIJtdPk2Yg1aKhPTQ5HYg+sP/rjts="}`
	session2 := `{"session_id":"d8b8304f5f339a818e127aca2dfab742","user_return_to":"/","_csrf_token":"atHkScP0CWcPrcxIJtdPk2Yg1aKhPTQ5HYg+sP/rjts="}`
	expectFailure(session, t)
	expectFailure(session2, t)
}

func expectFailure(session string, t *testing.T) {
	uid, err := ExtractUserId(session)
	if err == nil {
		t.Log("expected error to be returned")
		t.FailNow()
	}
	if err.Error() != WARDEN_FORMAT_ERROR {
		t.Logf("expected to fail with %s found %v", WARDEN_FORMAT_ERROR, err.Error())
		t.FailNow()
	}
	if uid != -1 {
		t.Logf("UserID in case of error should be -1 found %d\n", uid)
		t.FailNow()
	}
}
