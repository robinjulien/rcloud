import { Component, Input, OnInit, Output, EventEmitter } from '@angular/core';
import { FilemanagerService, FmFile, copycut } from 'src/app/services/filemanager.service';

@Component({
	selector: 'app-menu-action',
	templateUrl: './menu-action.component.html',
	styleUrls: ['./menu-action.component.css']
})
export class MenuActionComponent implements OnInit {
	@Input() path: string | undefined
	@Input() file: FmFile | undefined
	@Output() closeEvent: EventEmitter<null> = new EventEmitter<null>()
	@Output() refreshEvent: EventEmitter<null> = new EventEmitter<null>()
	@Output() copycutEvent: EventEmitter<copycut> = new EventEmitter<copycut>()
	editorMode: boolean = false
	content: string = ""

	constructor(private fm: FilemanagerService) { }

	ngOnInit(): void {
	}

	close(): void {
		this.closeEvent.next(null)
	}

	refresh(): void {
		this.refreshEvent.next(null)
	}

	rm(): void {
		if (window.confirm("Are you sure you want to delete " + this.file?.name + " ?")) {
			this.fm.rm(this.path + "/" + this.file?.name).subscribe(res => {
				if (res.success) {
					this.refresh()
				} else {
					window.alert(res.errorMessage)
				}
			})

			this.close()
		}
	}

	rename(): void {
		let newname = window.prompt("New name :", this.file?.name)

		if (newname != undefined && newname != "") {
			this.fm.mv(this.path + "/" + this.file?.name, this.path + "/" + newname).subscribe(res => {
				if (res.success) {
					this.refresh()
				} else {
					window.alert(res.errorMessage)
				}
			})

			this.close()
		}
	}

	download(): void {
		let downloadLink = document.createElement("a");
		downloadLink.href = "/api/fm/download?path=" + encodeURIComponent(this.path + "/" + this.file?.name);
		downloadLink.download = "" + this.file?.name;

		downloadLink.click()

		this.close()
	}

	cut() {
		this.copycutEvent.next({
			path: this.path + "",
			name: this.file?.name + "",
			operationIsCut: true // true is cut
		})

		this.close()
	}

	copy() {
		this.copycutEvent.next({
			path: this.path + "",
			name: this.file?.name + "",
			operationIsCut: false // false is copy
		})

		this.close()
	}

	openEditor() {
		this.fm.cat(this.path + "/" + this.file?.name).subscribe(res => {
			if (res.success) {
				this.content = res.content
			} else {
				window.alert(res.errorMessage)
			}
		})
		this.editorMode = true
	}

	closeEditor() {
		this.content = ""
		this.editorMode = false
	}

	save() {
		console.log(this.content)
		this.fm.echo(this.path + "/" + this.file?.name, this.content).subscribe(res => {
			if (res.success) {
				window.alert("Sucessfully saved.")
			} else {
				window.alert(res.errorMessage)
			}
		})
	}
}
