import { AfterViewInit, Component, ElementRef, OnInit, ViewChild } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { AuthService } from 'src/app/services/auth.service';
import { FilemanagerService, FmFile } from 'src/app/services/filemanager.service';

import '@grafikart/drop-files-element';
/*
<input *ngIf="render"
	  type="file"
	  multiple
	  name="files[]"
	  label="Drop files here or click to upload."
	  help="Upload files here and they won't be sent immediately"
	  is="drop-files"
/>
*/
@Component({
	selector: 'app-dashboard-view',
	templateUrl: './dashboard-view.component.html',
	styleUrls: ['./dashboard-view.component.css']
})
export class DashboardViewComponent implements OnInit, AfterViewInit {
	path: string = "."
	files: FmFile[] = []
	render: boolean = false
	@ViewChild("render") r!: ElementRef<HTMLDivElement>

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
		let upload = document.createElement("input")
		upload.type = "file"
		upload.multiple = true
		upload.name = "files[]"
		upload.setAttribute("is", "drop-files")
		this.r.nativeElement.innerHTML = '<input type="file" multiple is="drop-files"  >'
	}

	isRoot() {
		return [".", "/", "./", ""].includes(this.path)
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
}
