import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { HttpClientModule } from '@angular/common/http'

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { LoginViewComponent } from './pages/login/login-view/login-view.component';
import { DashboardViewComponent } from './pages/dashboard/dashboard-view/dashboard-view.component';
import { NavbarComponent } from './elements/navbar/navbar.component';
import { FormsModule } from '@angular/forms';
import { AuthService } from './services/auth.service';
import { LogoutViewComponent } from './pages/logout/logout-view/logout-view.component';

@NgModule({
  declarations: [
    AppComponent,
    LoginViewComponent,
    DashboardViewComponent,
    NavbarComponent,
    LogoutViewComponent,
  ],
  imports: [
    BrowserModule,
	AppRoutingModule,
	HttpClientModule,
	FormsModule
  ],
  providers: [
	  AuthService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
