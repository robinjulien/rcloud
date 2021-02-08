import { Component, OnInit } from '@angular/core';
import { NgForm } from '@angular/forms';
import { Router } from '@angular/router';
import { AuthService } from 'src/app/services/auth.service';

@Component({
	selector: 'app-login-view',
	templateUrl: './login-view.component.html',
	styleUrls: ['./login-view.component.css']
})
export class LoginViewComponent implements OnInit {
	wrongidpwd: boolean = false

	constructor(private auth: AuthService, private router: Router) { }

	ngOnInit(): void {
	}

	onSubmit(f: NgForm) {
		let id = f.value["id"]
		let password = f.value["password"]

		if (id == undefined || id == "" || password == undefined || password == "") {
			return
		}

		this.auth.attemptLogin(id, password).subscribe(res => {
			if (res.success) {
				this.auth.loadStatus()
				this.router.navigate(["/"])
			} else {
				this.wrongidpwd = true
			}
		})
	}
}
