package main

import "net/http"

func (app *application) enableCORS(next http.Handler) http.Handler {
	// CORS (Corss-Origin Resource Sharing)
	// : HTTPヘッダに"Access-Control-Allow-Origin"を付加して、あるオリジンのアプリケーションに、別のオリジンのアプリケーションへのアクセス権を与えること。
	//   今回の場合、ポート3000で動くフロントアプリが、ポート8080で動くバックエンドアプリへのアクセスを可能にしている。

	// オリジン：プロトコル、ドメイン、ポート番号

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// *: Allow all requests
		w.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r)
	})
}
