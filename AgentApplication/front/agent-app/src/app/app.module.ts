import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import {MatTableModule} from '@angular/material/table';
import { MatSliderModule } from '@angular/material/slider';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule  } from '@angular/material/icon';
import { MatGridListModule} from '@angular/material/grid-list';
import { MatCardModule} from '@angular/material/card';
import { MatBadgeModule} from '@angular/material/badge';
import { MatDialogModule} from '@angular/material/dialog';
import { MatFormFieldModule} from '@angular/material/form-field';
import { MatInputModule} from '@angular/material/input';
import { MatCheckboxModule} from '@angular/material/checkbox';
import { MatStepperModule} from '@angular/material/stepper';
import { MatRadioModule} from '@angular/material/radio';
import { MatDividerModule} from '@angular/material/divider';
import { MatListModule} from '@angular/material/list';
import { MatDatepickerModule} from '@angular/material/datepicker';
import { MatNativeDateModule } from '@angular/material/core';
import { MatSelectModule} from '@angular/material/select';
import {MatSnackBarModule} from '@angular/material/snack-bar';
import {MatTabsModule} from '@angular/material/tabs';
import {MatMenuModule} from '@angular/material/menu';
import {MatPaginatorModule} from '@angular/material/paginator';
import { MatCarouselModule } from 'ng-mat-carousel';
import {MatSortModule} from '@angular/material/sort';
import {Sort} from '@angular/material/sort';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { JwtHelperService, JWT_OPTIONS } from '@auth0/angular-jwt';
import { InterceptorService } from './service/interceptor.service';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
import { LoginPageComponent } from './component/login-page/login-page.component';
import { RegistrationPageComponent } from './component/registration-page/registration-page.component';
import { UserService } from './service/user.service';
import { LandingPageComponent } from './component/landing-page/landing-page.component';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { HeaderComponent } from './component/header/header.component';
import { FooterComponent } from './component/footer/footer.component';
import { CompaniesComponent } from './component/companies/companies.component';
import { UserProfileComponent } from './component/user-profile/user-profile.component';
import { CreateCompanyComponent } from './component/create-company/create-company.component';
import { ClientCompaniesComponent } from './component/client-companies/client-companies.component';
import { ForumComponent } from './component/forum/forum.component';
import { CompanyService } from './service/company.service';

const MaterialComponents = [
  MatSliderModule,
  MatCarouselModule,
  MatToolbarModule,
  MatButtonModule,
  MatIconModule,
  MatGridListModule,
  MatCardModule,
  MatBadgeModule,
  MatDialogModule,
  MatFormFieldModule,
  MatInputModule,
  MatCheckboxModule,
  MatStepperModule,
  MatRadioModule,
  MatDividerModule,
  MatListModule,
  MatDatepickerModule,
  MatNativeDateModule,
  MatSelectModule,
  MatSnackBarModule,
  MatTabsModule,
  MatMenuModule,
  MatTableModule,
  MatPaginatorModule,
  MatSortModule
];

@NgModule({
  declarations: [
    AppComponent,
    LoginPageComponent,
    RegistrationPageComponent,
    LandingPageComponent,
    HeaderComponent,
    FooterComponent,
    CompaniesComponent,
    UserProfileComponent,
    CreateCompanyComponent,
    ClientCompaniesComponent,
    ForumComponent
  ],
  imports: [
    MaterialComponents,
    BrowserModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    BrowserModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    MaterialComponents,
    HttpClientModule,
    FormsModule,
    ReactiveFormsModule,
    CommonModule
  ],
  providers: [UserService, CompanyService, { provide: JWT_OPTIONS, useValue: JWT_OPTIONS }, JwtHelperService,
    { provide: HTTP_INTERCEPTORS, useClass: InterceptorService, multi: true }],
  bootstrap: [AppComponent]
})
export class AppModule { }
