import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';

export interface FmFile {
	name: string
	size: number
	isDir: boolean
}

interface responseLs {
	success: boolean
	errorMessage: string
	dir: FmFile[]
}

@Injectable({
	providedIn: 'root'
})
export class FilemanagerService {
	constructor(private http: HttpClient) { }

	ls(path: string) {
		let params = new HttpParams().set("path", path)
		return this.http.get<responseLs>("/api/fm/ls", {params: params})
	}
}
