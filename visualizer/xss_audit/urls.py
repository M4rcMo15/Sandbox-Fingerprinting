from django.urls import path
from . import views

app_name = 'xss_audit'

urlpatterns = [
    # Callback público (sin auth) para recibir XSS
    path('xss-callback', views.xss_callback, name='xss_callback'),
    
    # Dashboard y vistas (con auth)
    path('dashboard/', views.xss_dashboard, name='dashboard'),
    path('hit/<int:hit_id>/', views.xss_hit_detail, name='hit_detail'),
    path('payloads/', views.xss_payloads_list, name='payloads_list'),
    path('sandbox/<int:sandbox_id>/', views.sandbox_detail, name='sandbox_detail'),
]
