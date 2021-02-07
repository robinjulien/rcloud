import { stringify } from '@angular/compiler/src/util';
import { Pipe, PipeTransform } from '@angular/core';
import { FmFile } from '../services/filemanager.service';

@Pipe({
	name: 'sortfilesalphabetical'
})
export class SortfilesalphabeticalPipe implements PipeTransform {

	transform(items: FmFile[], ...args: unknown[]): FmFile[] {
		return items.sort((a, b) => { return a.name.toUpperCase().localeCompare(b.name.toUpperCase()) });
	}

}
