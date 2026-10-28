import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterLink],
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.css']
})
export class RegisterComponent {
  private authService = inject(AuthService);
  private router = inject(Router);

  username = '';
  password = '';
  confirmPassword = '';
  errorMessage = '';
  successMessage = '';
  isLoading = false;

  onSubmit(): void {
    this.errorMessage = '';
    this.successMessage = '';

    if (!this.username.trim() || !this.password || !this.confirmPassword) {
      this.errorMessage = 'กรุณากรอกข้อมูลให้ครบทุกช่อง';
      return;
    }

    if (this.password !== this.confirmPassword) {
      this.errorMessage = 'รหัสผ่านและยืนยันรหัสผ่านไม่ตรงกัน';
      return;
    }

    this.isLoading = true;

    this.authService.register({
      username: this.username.trim(),
      password: this.password,
      confirm_password: this.confirmPassword
    }).subscribe({
      next: () => {
        this.isLoading = false;
        this.successMessage = 'สมัครสมาชิกสำเร็จ! กำลังนำท่านไปยังหน้าลงชื่อเข้าใช้งาน...';

        setTimeout(() => {
          this.router.navigate(['/login']);
        }, 1500);
      },
      error: (err) => {
        this.isLoading = false;
        if (err.error?.error) {
          this.errorMessage = err.error.error;
        } else {
          this.errorMessage = 'เกิดข้อผิดพลาดในการสมัครสมาชิก กรุณาลองใหม่อีกครั้ง';
        }
      }
    });
  }
}
