import { Component, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { AuthService, ProfileResponse } from '../../services/auth.service';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {
  private authService = inject(AuthService);
  private router = inject(Router);

  profile: ProfileResponse | null = null;
  isLoading = true;
  errorMessage = '';

  ngOnInit(): void {
    this.loadUserProfile();
  }

  loadUserProfile(): void {
    const token = this.authService.getToken();
    if (!token) {
      this.handleAuthFailure('ไม่พบ Token ยืนยันตัวตน กรุณาเข้าสู่ระบบใหม่');
      return;
    }

    this.isLoading = true;

    this.authService.getProfile().subscribe({
      next: (data) => {
        this.profile = data;
        this.isLoading = false;
      },
      error: (err) => {
        console.error('Profile fetch failed:', err);
        this.handleAuthFailure('เซสชันหมดอายุหรือไม่ถูกต้อง กรุณาเข้าสู่ระบบใหม่อีกครั้ง');
      }
    });
  }

  private handleAuthFailure(reason: string): void {
    this.isLoading = false;
    this.authService.clearToken();
    alert(reason);
    this.router.navigate(['/login']);
  }

  logout(): void {
    this.authService.clearToken();
    this.router.navigate(['/login']);
  }
}
