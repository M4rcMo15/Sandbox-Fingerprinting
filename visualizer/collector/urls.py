from django.urls import path
from . import views

urlpatterns = [
    path('', views.index, name='index'),
    path('statistics/', views.statistics, name='statistics'),
    path('api/collect/', views.collect_data, name='collect_data'),
    path('api/collect', views.collect_data, name='collect_data_no_slash'),  # Sin barra también
    path('execution/<uuid:guid>/', views.execution_detail, name='execution_detail'),
]
