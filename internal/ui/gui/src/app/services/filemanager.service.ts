import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

export interface FmFile {
	name: string
	size: number
	isDir: boolean
}

export interface responseLs {
	success: boolean
	errorMessage: string
	dir: FmFile[]
}

export interface BaseResponse {
	success: boolean
	errorMessage: string
}

@Injectable({
	providedIn: 'root'
})
export class FilemanagerService {
	constructor(private http: HttpClient) { }

	ls(path: string): Observable<responseLs> {
		let params = new HttpParams().set("path", path)
		return this.http.get<responseLs>("/api/fm/ls", {params: params})
	}

	mkdir(path: string): Observable<BaseResponse> {
		let fd = new FormData()
		fd.append("path", path)
		return this.http.post<BaseResponse>("/api/fm/mkdir", fd)
	}

	touch(path: string): Observable<BaseResponse> {
		let fd = new FormData()
		fd.append("path", path)
		return this.http.post<BaseResponse>("/api/fm/touch", fd)
	}

	upload(): void {
		// upload isn't managed by angular httpclient, so its definition is in the dashboard logic
	}
}
