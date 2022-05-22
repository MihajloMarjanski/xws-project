import { Injectable } from '@angular/core';
import { MatSnackBar } from '@angular/material/snack-bar';
import { ActivatedRouteSnapshot, CanActivate, RouterStateSnapshot, UrlTree, Router } from '@angular/router';
import { Observable } from 'rxjs';
import { AuthService } from "./auth.service";

@Injectable({
  providedIn: 'root'
})
export class AuthGuard implements CanActivate {

  constructor(private authService: AuthService, private router: Router, private _snackBar: MatSnackBar) { }

  canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): boolean | Promise<boolean> {
        var isAuthenticated = this.authService.getAuthStatus();
        var hasExpired = this.authService.hasExpired();

        if (!isAuthenticated || hasExpired ) {          
          this.router.navigate(['/login']).then(_ =>
            this._snackBar.open('You have to login first', 'Close', {duration: 5000}))
          return false;
        }

        if (!route.data.role.includes(localStorage.getItem('roles'))) {          
          this.router.navigate(['/']).then(_ =>
                this._snackBar.open('You are not authorized for this funcionality', 'Close', {duration: 5000}))
          return false;
        }          

        return true;
    }
  
}
