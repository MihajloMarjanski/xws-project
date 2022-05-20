import { Component, OnInit } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Company } from 'src/app/model/company';
import { UserService } from 'src/app/service/user.service';

@Component({
  selector: 'app-companies',
  templateUrl: './companies.component.html',
  styleUrls: ['./companies.component.css']
})
export class CompaniesComponent implements OnInit {

  companies: Company[] = []
  displayedColumns: string[] = ['name', 'info', 'owner', 'approve'];
  errorMessage : string  = '';

  constructor(private _userService: UserService, private _snackBar: MatSnackBar) { }

  ngOnInit(): void {
    this.getCompanies();
  }
  

  getCompanies() {
    this._userService.getCompanies()
        .subscribe(data =>  {
                    this.companies = data   
                    },
                   error => this.errorMessage = <any>error);   
  }

  approve(id: number) {
    this._userService.approveCompany(id)
        .subscribe(data =>  {  
                    },
                   error => this.errorMessage = <any>error);
  }

}
