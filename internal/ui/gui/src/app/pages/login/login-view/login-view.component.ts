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

	constructor(private auth: AuthService, private router: Router) { }

	ngOnInit(): void {
	}

	onSubmit(f: NgForm) {
		this.auth.attemptLogin(f.value["id"], f.value["password"]).subscribe(res => {
			if (res.success) {
				this.auth.loadStatus()
				this.router.navigate(["/"])
			} else {
				// Afficher message d'erreur
			}
		})
	}
}
