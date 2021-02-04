import { Pipe, PipeTransform } from '@angular/core';
import { FmFile } from '../services/filemanager.service';

@Pipe({
	name: 'isnotdir'
})
export class IsnotdirPipe implements PipeTransform {

	transform(items: FmFile[], ...args: unknown[]): FmFile[] {
		return items.filter(item => !item.isDir);
	}

}
