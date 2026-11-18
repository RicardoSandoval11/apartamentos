using Microsoft.VisualBasic;
using scheduler_service.DataAccess.Config;
using System;

var builder = WebApplication.CreateBuilder(args);

// Add services to the container.
// Learn more about configuring OpenAPI at https://aka.ms/aspnet/openapi
builder.Services.AddOpenApi();
builder.Services.AddControllers();

string? connectionString = builder.Configuration["DatabaseSettings:ConnectionString"];

if (string.IsNullOrEmpty(connectionString))
{
    Environment.Exit(1);
}

builder.Services.AddDatabase(connectionString);
builder.Services.ApplyMigrations();

var app = builder.Build();

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.MapOpenApi();
}

app.UseHttpsRedirection();
app.MapControllers();

await app.RunAsync();
