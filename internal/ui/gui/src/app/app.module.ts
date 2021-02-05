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
import { FilemanagerService } from './services/filemanager.service';
import { HumanreadablesizePipe } from './pipes/humanreadablesize.pipe';
import { ExttoimgPipe } from './pipes/exttoimg.pipe';
import { IsdirPipe } from './pipes/isdir.pipe';
import { IsnotdirPipe } from './pipes/isnotdir.pipe';
import { SortfoldersfilesPipe } from './pipes/sortfoldersfiles.pipe';
import { SortfilesalphabeticalPipe } from './pipes/sortfilesalphabetical.pipe';
import { MenuActionComponent } from './pages/dashboard/menu-action/menu-action.component';

@NgModule({
  declarations: [
    AppComponent,
    LoginViewComponent,
    DashboardViewComponent,
    NavbarComponent,
    LogoutViewComponent,
    HumanreadablesizePipe,
    ExttoimgPipe,
    IsdirPipe,
    IsnotdirPipe,
    SortfoldersfilesPipe,
    SortfilesalphabeticalPipe,
    MenuActionComponent,
  ],
  imports: [
    BrowserModule,
	AppRoutingModule,
	HttpClientModule,
	FormsModule
  ],
  providers: [
	  AuthService,
	  FilemanagerService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
