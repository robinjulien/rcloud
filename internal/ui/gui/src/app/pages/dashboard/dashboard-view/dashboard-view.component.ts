import { AfterViewInit, Component, ElementRef, OnInit, ViewChild } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { AuthService } from 'src/app/services/auth.service';
import { FilemanagerService, FmFile, responseLs } from 'src/app/services/filemanager.service';

function fileNameInLs(name: string, dir: FmFile[]): boolean {
	for (let file of dir) {
		if (name == file.name) {
			return true
		}
	}
	return false
}

@Component({
	selector: 'app-dashboard-view',
	templateUrl: './dashboard-view.component.html',
	styleUrls: ['./dashboard-view.component.css']
})
export class DashboardViewComponent implements OnInit, AfterViewInit {
	path: string = "."
	files: FmFile[] = []
	render: boolean = false
	uploadProgress: number = 0
	selectedFile: FmFile | undefined

	constructor(private auth: AuthService, private fm: FilemanagerService, private router: Router, private route: ActivatedRoute) { }

	ngOnInit(): void {
		this.auth.amILoggedIn().subscribe(res => {
			if (res.loggedIn) {
				/*this.route.queryParams.subscribe(p => {
					if (this.route.snapshot.queryParams.path != undefined && this.route.snapshot.queryParams.path != null) {
						this.path = this.route.snapshot.queryParams.path
					}
					this.ls()
				})*/

				this.ls()
			} else {
				this.router.navigate(["/login"])
			}
		})
	}

	ngAfterViewInit() {
	}

	isRoot() {
		return [".", "/", "./", ""].includes(this.path)
	}

	clickRow(e: MouseEvent, i: number) {
		if (!(e.target instanceof HTMLButtonElement)) {
			if (this.files[i].isDir) {
				this.navigate(i)
			} else {
				this.download(i)
			}
		}
	}

	showActions(i: number) {
		this.selectedFile = undefined
		this.selectedFile = this.files[i]
	}

	navigate(i: number) {
		if (i < 0) {
			let split = this.path.split("/")
			this.path = split.length > 1 ? split.slice(0, split.length - 1).join("/") : this.path
		} else {
			this.path = this.path + "/" + this.files[i].name
		}

		/*
		this.router.navigate([], {
			relativeTo: this.route,
			queryParams: { path: this.path },
		});
		
		if (history.pushState) {
			var newurl = window.location.protocol + "//" + window.location.host + window.location.pathname + '?myNewUrlQuery=1';
			window.history.pushState({path:newurl},'',newurl);
		}*/

		this.ls()
	}

	download(i: number) {
		let downloadLink = document.createElement("a");
		downloadLink.href = "/api/fm/download?path=" + encodeURIComponent(this.path + "/" + this.files[i].name);
		downloadLink.download = this.files[i].name;

		downloadLink.click()
	}

	ls(): void {
		this.fm.ls(this.path).subscribe(res => {
			this.files = res.dir
		})
	}

	mkdir() {
		let foldername = window.prompt("Name of the new folder :")

		if (foldername == null || foldername == "") {
			return
		}

		this.fm.mkdir(this.path + "/" + foldername).subscribe(res => {
			if (!res.success) {
				window.alert(res.errorMessage)
			} else {
				this.ls()
			}
		})
	}

	touch() {
		let filename = window.prompt("Name of the new file :")

		if (filename == null || filename == "") {
			return
		}

		this.fm.touch(this.path + "/" + filename).subscribe(res => {
			if (!res.success) {
				window.alert(res.errorMessage)
			} else {
				this.ls()
			}
		})
	}

	upload(iff: HTMLInputElement): void {
		let params = new URLSearchParams()
		params.append("path", this.path)

		fetch("/api/fm/ls?" + params.toString(), {
			method: "GET",
		}).then(res => res.json()).then(res => {
			let resls = res as responseLs

			if (resls.success) {
				let req = new XMLHttpRequest();
				req.open("post", "/api/fm/upload", true);

				let fd = new FormData()

				let length = iff != null ? (iff.files != null ? iff.files.length : 0) : 0

				if (length == 0) {
					return
				}

				for (let i = 0; i < length; i++) {
					let f = iff != null ? (iff.files != null ? iff.files[i] : "") : ""

					let file = f as File
					if (fileNameInLs(file.name, resls.dir)) {
						if (!window.confirm("File " + file.name + " already exists in this directory. Do you want to overwrite it ?")) {
							continue // If response is false, go to next iteration without adding the file
						}
					}

					fd.append("files[]", file)
				}

				fd.append("path", this.path)

				req.upload.onprogress = e => {
					this.uploadProgress = e.loaded / e.total * 100
				}

				// Automatic multipart/form-data
				req.send(fd);
			}
		})
	}
}
