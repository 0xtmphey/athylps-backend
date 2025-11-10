package hooks

import (
	"net/http"
)

func HandleDonationAlertsWebhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "text/html")
		w.Write([]byte("<html>ok</html>"))
	}
}
