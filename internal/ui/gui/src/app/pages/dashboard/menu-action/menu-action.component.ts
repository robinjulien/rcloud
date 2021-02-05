import { Component, Input, OnInit, Output, EventEmitter } from '@angular/core';
import { FmFile } from 'src/app/services/filemanager.service';

@Component({
	selector: 'app-menu-action',
	templateUrl: './menu-action.component.html',
	styleUrls: ['./menu-action.component.css']
})
export class MenuActionComponent implements OnInit {
	@Input() file: FmFile | undefined
	@Output() close: EventEmitter<null> = new EventEmitter<null>()

	constructor() { }

	ngOnInit(): void {
	}

	closeActions() {
		this.close.next(null)
	}
}
