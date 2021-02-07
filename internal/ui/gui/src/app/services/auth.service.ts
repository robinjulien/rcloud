import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { BaseResponse } from './filemanager.service';

interface responseLogin {
	success: boolean
}

interface responseWhoAmI {
	loggedIn: boolean
	id: string
	admin: boolean
}

interface responseAmILoggedIn {
	loggedIn: boolean
}

export interface PublicUser {
	id: string
	admin: boolean
}

@Injectable({
	providedIn: 'root'
})
export class AuthService {
	private loggedIn: boolean = false
	private user!: PublicUser

	constructor(private http: HttpClient) {
		this.loadStatus()
	}

	loadStatus() {
		this.http.get<responseWhoAmI>("/api/auth/whoami").subscribe(res => {
			if (res.loggedIn) {
				this.loggedIn = true
				this.user = { id: res.id, admin: res.admin }
			} else {
				this.loggedIn = false
				this.user = { id: "", admin: false }
			}
		})
	}

	amILoggedIn(): Observable<responseAmILoggedIn> {
		return this.http.get<responseAmILoggedIn>("/api/auth/amiloggedin")
	}

	whoAmI() {
		return this.http.get<responseWhoAmI>("/api/auth/whoami")
	}

	isLoggedIn(): boolean {
		return this.loggedIn
	}

	getUser(): PublicUser {
		return this.user
	}

	attemptLogin(id: string, password: string): Observable<responseLogin> {
		let body = new FormData()
		body.append("id", id)
		body.append("password", password)
		return this.http.post<responseLogin>("/api/auth/login", body)
	}

	logout(): Observable<any> {
		return this.http.post("/api/auth/logout", null)
	}

	changeMyPassword(oldpwd: string, newpwd: string) {
		let fd = new FormData()
		fd.append("oldpassword", oldpwd)
		fd.append("newpassword", newpwd)
		return this.http.post<BaseResponse>("/api/auth/change-password", fd)
	}
}
