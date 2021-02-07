import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { PublicUser } from './auth.service';
import { BaseResponse } from './filemanager.service';

interface responseListUsers {
	success: boolean
	errorMessage: string
	users: PublicUser[]
}

@Injectable({
	providedIn: 'root'
})
export class UsersService {

	constructor(private http: HttpClient) { }

	listUsers(): Observable<responseListUsers> {
		return this.http.get<responseListUsers>("/api/users/list")
	}

	addUser(id: string, pwd: string, admin: boolean): Observable<BaseResponse> {
		let fd = new FormData()
		fd.append("id", id)
		fd.append("password", pwd)
		fd.append("admin", admin ? "true" : "false")
		return this.http.post<BaseResponse>("/api/users/add", fd)
	}

	delUser(id: string): Observable<BaseResponse> {
		let fd = new FormData()
		fd.append("id", id)
		return this.http.post<BaseResponse>("/api/users/del", fd)
	}

	editUserPassword(id: string, password: string): Observable<BaseResponse> {
		let fd = new FormData()
		fd.append("id", id)
		fd.append("password", password)
		return this.http.post<BaseResponse>("/api/users/edit", fd)
	}

	editUserAdmin(id: string, admin: boolean): Observable<BaseResponse> {
		let fd = new FormData()
		fd.append("id", id)
		fd.append("admin", admin ? "true" : "false")
		return this.http.post<BaseResponse>("/api/users/edit", fd)
	}
}
