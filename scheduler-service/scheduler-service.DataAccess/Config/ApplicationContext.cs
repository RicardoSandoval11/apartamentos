using Microsoft.EntityFrameworkCore;

namespace scheduler_service.DataAccess.Config;

public class ApplicationContext(DbContextOptions<ApplicationContext> options) : DbContext(options)
{
    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.ApplyConfiguration(new ReservationConfiguration());
        modelBuilder.ApplyConfiguration(new TimeslotConfiguration());
        base.OnModelCreating(modelBuilder);
    }
}