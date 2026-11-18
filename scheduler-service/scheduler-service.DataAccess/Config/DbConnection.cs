using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.DependencyInjection;

namespace scheduler_service.DataAccess.Config;

public static class DbConnection
{
    public static IServiceCollection AddDatabase(this IServiceCollection services, string connectionString)
    {
        services.AddDbContext<ApplicationContext>(options => options.UseNpgsql(connectionString));
        return services;
    }

    public static IServiceCollection ApplyMigrations(this IServiceCollection services)
    {
        using (var serviceProvider = services.BuildServiceProvider())
        {
            using (var scope = serviceProvider.CreateScope())
            {
                using (var dbContext = scope.ServiceProvider.GetRequiredService<ApplicationContext>())
                {
                    var pendingMigrations = dbContext.Database.GetPendingMigrations();

                    if (pendingMigrations.Any())
                    {
                        dbContext.Database.Migrate();
                    }
                }
            }
        }
        return services;
    }
}