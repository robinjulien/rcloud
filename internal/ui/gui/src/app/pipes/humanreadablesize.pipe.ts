import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
	name: 'humanreadablesize'
})
export class HumanreadablesizePipe implements PipeTransform {

	transform(value: number, ...args: unknown[]): string {
		const values = ["B", "KB", "MB", "GB", "TB"]
		let i = 0
		let a = value

		for (; a >= 1000; i++) {
			a /= 1000
		}

		return Math.round(a) + " " + values[i];
	}

}
