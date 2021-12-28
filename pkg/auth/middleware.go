package auth

import "net/http"

func Middleware(authSvc AuthService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u, err := authSvc.UserFrom(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}

			ctx := r.Context()
			ctx = ToContext(ctx, u)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
