package pgconnectionstring

import "testing"

func TestConnectionStringSucces(t *testing.T) {
	tests := []struct {
		connstr string
		err     error
		values  map[string]string
	}{
		{
			connstr: `user=user dbname=db password=pwd`,
			err:     nil,
			values: map[string]string{
				"user":     "user",
				"dbname":   "db",
				"password": "pwd",
			},
		},
		{
			connstr: `user='user' dbname='db' password='pwd'`,
			err:     nil,
			values: map[string]string{
				"user":     "user",
				"dbname":   "db",
				"password": "pwd",
			},
		},
		{
			connstr: `user user`,
			err:     ErrMissingEqualSign,
			values: map[string]string{
				"user": "user",
			},
		},
		{
			connstr: `user=user\`,
			err:     ErrMissingCharacterAfterBackslash,
			values: map[string]string{
				"user": "user",
			},
		},
		{
			connstr: `user='user`,
			err:     ErrUnterminatedQuote,
			values: map[string]string{
				"user": "user",
			},
		},
	}

	for _, v := range tests {
		values, err := Parse(v.connstr)
		if err != v.err {
			t.Errorf(`Errors are not equal in "%s"`, v.connstr)
			continue
		}
		if err != nil {
			continue
		}
		for i, o := range v.values {
			if value, exists := values[i]; exists {
				if value != o {
					t.Errorf(`Value "%s" is not equal to value "%s" in result values of connection string "%s"`, value, i, v.connstr)
					continue
				}
			} else {
				t.Errorf(`Key "%s" is not exists in result values of connection string "%s"`, i, v.connstr)
				continue
			}
		}
	}
}
