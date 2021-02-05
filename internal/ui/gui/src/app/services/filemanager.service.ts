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

export interface copycut {
	path: string
	name: string
	operationIsCut: boolean
}

export interface responseCat{
	success: boolean
	errorMessage: string
	content: string
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

	rm(path: string): Observable<BaseResponse> {
		let fd = new FormData()
		fd.append("path", path)
		return this.http.post<BaseResponse>("/api/fm/rm", fd)
	}

	mv(src: string, dest: string): Observable<BaseResponse> {
		let fd = new FormData()
		fd.append("src", src)
		fd.append("dest", dest)
		return this.http.post<BaseResponse>("/api/fm/mv", fd)
	}

	copy(src: string, dest: string) {
		let fd = new FormData()
		fd.append("src", src)
		fd.append("dest", dest)
		return this.http.post<BaseResponse>("/api/fm/cp", fd)
	}

	cat(path: string): Observable<responseCat> {
		let params = new HttpParams().set("path", path)
		return this.http.get<responseCat>("/api/fm/cat", {params: params})
	}

	echo(path: string, content: string): Observable<BaseResponse> {
		let fd = new FormData()
		fd.append("path", path)
		fd.append("content", content)
		return this.http.post<BaseResponse>("/api/fm/echo", fd)
	}
}
