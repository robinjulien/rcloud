import { Pipe, PipeTransform } from '@angular/core';
import { FmFile } from '../services/filemanager.service';

@Pipe({
  name: 'sortfoldersfiles'
})
export class SortfoldersfilesPipe implements PipeTransform {

  transform(items: FmFile[], ...args: unknown[]): FmFile[] {
    return items.sort((a, b) => {
		if (a.isDir == b.isDir) {
			return 0
		} else if (a.isDir) {
			return -1
		} else {
			return 1
		}
	});
  }

}
