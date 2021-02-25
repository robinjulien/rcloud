import { AfterViewInit, Component, ElementRef, OnInit, ViewChild } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { AuthService } from 'src/app/services/auth.service';
import { BaseResponse, FilemanagerService, FmFile, responseLs, copycut } from 'src/app/services/filemanager.service';

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
	lsLoaded: boolean = false
	files: FmFile[] = []
	render: boolean = false
	uploadProgress: number = 0
	selectedFile: FmFile | undefined
	copycut: copycut = {path: "", name: "", operationIsCut: true}

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

	isAdmin(): boolean {
		return this.auth.isLoggedIn() && this.auth.getUser().admin
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

		this.files = []
		this.ls()
	}

	download(i: number) {
		let downloadLink = document.createElement("a");
		downloadLink.href = "/api/fm/download?path=" + encodeURIComponent(this.path + "/" + this.files[i].name);
		downloadLink.download = this.files[i].name;

		downloadLink.click()
	}

	ls(): void {
		this.lsLoaded = false
		this.fm.ls(this.path).subscribe(res => {
			this.files = res.dir
			this.lsLoaded = true
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

	paste(): void {
		if (this.copycut.name != "") {
			if (this.copycut.operationIsCut) {
				this.fm.mv(this.copycut.path + "/" + this.copycut.name, this.path + "/" + this.copycut.name).subscribe(res => {
					if (res.success) {
						this.ls()
					} else {
						window.alert(res.errorMessage)
					}
				})
			} else {
				this.fm.copy(this.copycut.path + "/" + this.copycut.name, this.path + "/" + this.copycut.name).subscribe(res => {
					if (res.success) {
						this.ls()
					} else {
						window.alert(res.errorMessage)
					}
				})
			}
		}
	}

	upload(iff: HTMLInputElement): void { // iff = input files form
		let params = new URLSearchParams()
		params.append("path", this.path)

		fetch("/api/fm/ls?" + params.toString(), { // get list of files to check if it will overwrite
			method: "GET",
		}).then(res => res.json()).then(res => {
			let resls = res as responseLs

			if (resls.success) {
				let req = new XMLHttpRequest();
				req.open("post", "/api/fm/upload", true);

				let fd = new FormData()

				let length = iff != null ? (iff.files != null ? iff.files.length : 0) : 0 // check if iff or iff.files is needed to get length ?

				if (length == 0) {
					return
				}

				// Build the formdata with all files
				for (let i = 0; i < length; i++) {
					let f = iff != null ? (iff.files != null ? iff.files[i] : "") : "" // check if iff or iff.files is needed to get file

					let file = f as File
					if (fileNameInLs(file.name, resls.dir)) { // for each overwrite, ask consent
						if (!window.confirm("File " + file.name + " already exists in this directory. Do you want to overwrite it ?")) {
							continue // If response is false, go to next iteration without adding the file
						}
					}

					fd.append("files[]", file)
				}

				fd.append("path", this.path)

				// Upload progress updates the progressBar
				req.upload.onprogress = e => {
					this.uploadProgress = e.loaded / e.total * 100
				}

				// On finish, progressbar disappear and if there is an error, it is shown
				req.onreadystatechange = e => {
					if (req.readyState == 4 && req.status == 200) {
						let json = JSON.parse(req.responseText) as BaseResponse

						if (!json.success) {
							window.alert(json.errorMessage)
						} else {
							this.ls()
						}

						this.uploadProgress = 0
					}
				}

				// Automatic multipart/form-data
				req.send(fd);
			}
		})
	}
}
