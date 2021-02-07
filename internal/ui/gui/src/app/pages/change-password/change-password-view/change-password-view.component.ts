import { Component, OnInit } from '@angular/core';
import { NgForm } from '@angular/forms';
import { Router } from '@angular/router';
import { AuthService } from 'src/app/services/auth.service';

@Component({
	selector: 'app-change-password-view',
	templateUrl: './change-password-view.component.html',
	styleUrls: ['./change-password-view.component.css']
})
export class ChangePasswordViewComponent implements OnInit {
	wrongpwd: boolean = false
	wrongconf: boolean = false

	constructor(private auth: AuthService, private router: Router) { }

	ngOnInit(): void {
	}

	onSubmit(f: NgForm) {
		this.wrongconf = false
		this.wrongpwd = false

		let oldpwd = f.value["oldpwd"]
		let newpwd = f.value["newpwd"]
		let newpwdc = f.value["newpwdc"]

		if (oldpwd == undefined || newpwd == undefined || newpwdc == undefined || oldpwd == "" || newpwd == "" || newpwdc == "") {
			return
		}

		if (newpwd != newpwdc) {
			this.wrongconf = true
			return
		}

		this.auth.changeMyPassword(oldpwd, newpwd).subscribe(res => {
			if (res.success) {
				window.alert("Done !")
				this.router.navigate(["/"])
			} else {
				this.wrongpwd = true
			}
		})
	}
}
