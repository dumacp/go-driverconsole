package restclient

import (
	"fmt"
	"net/http"

	"github.com/dumacp/go-driverconsole/internal/utils"
)

func GetData(c *http.Client, id, urlin, filterHttpQuery string) ([]byte, int, error) {
	if c == nil {
		return nil, 0, fmt.Errorf("client http empty")
	}

	url := fmt.Sprintf("%s%s", urlin, filterHttpQuery)
	resp, code, err := utils.Get(c, url, "", "", nil)
	if err != nil {
		return nil, code, err
	}

	return resp, code, nil
}
