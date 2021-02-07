import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService, PublicUser } from 'src/app/services/auth.service';
import { UsersService } from 'src/app/services/users.service';

@Component({
	selector: 'app-admin-view',
	templateUrl: './admin-view.component.html',
	styleUrls: ['./admin-view.component.css']
})
export class AdminViewComponent implements OnInit {
	users!: PublicUser[]

	constructor(private auth: AuthService, private usersService: UsersService, private router: Router) { }

	ngOnInit(): void {
		this.auth.whoAmI().subscribe(res => {
			if (res.loggedIn) {
				if (res.admin) {
					this.listUsers()
				} else {
					this.router.navigate(["/"])
				}
			} else {
				this.router.navigate(["/login"])
			}
		})
	}

	listUsers(): void {
		this.usersService.listUsers().subscribe(res => {
			if (res.success) {
				this.users = res.users
			} else {
				window.alert(res.errorMessage)
			}
		})
	}

	addUser(): void {
		let id = window.prompt("User id :")
		let pwd = window.prompt("User password :")

		if (id != undefined && pwd != undefined && id != "" && pwd != "") {
			this.usersService.addUser(id, pwd, false).subscribe(res => {
				if (res.success) {
					this.listUsers()
				} else {
					window.alert(res.errorMessage)
				}
			})
		}

	}

	deleteUser(i: number) {
		this.usersService.delUser(this.users[i].id).subscribe(res => {
			if (res.success) {
				this.listUsers()
			} else {
				window.alert(res.errorMessage)
			}
		})
	}

	promoteUser(i: number) {
		this.usersService.editUserAdmin(this.users[i].id, true).subscribe(res => {
			if (res.success) {
				this.listUsers()
			} else {
				window.alert(res.errorMessage)
			}
		})
	}

	demoteUser(i: number) {
		this.usersService.editUserAdmin(this.users[i].id, false).subscribe(res => {
			if (res.success) {
				this.listUsers()
			} else {
				window.alert(res.errorMessage)
			}
		})
	}

	changePasswordUser(i: number) {
		let new_pwd = window.prompt("New password for the user " + this.users[i].id + " :")

		if (new_pwd != undefined && new_pwd != "") {
			this.usersService.editUserPassword(this.users[i].id, new_pwd).subscribe(res => {
				if (res.success) {
					this.listUsers()
				} else {
					window.alert(res.errorMessage)
				}
			})
		}
	}
}
