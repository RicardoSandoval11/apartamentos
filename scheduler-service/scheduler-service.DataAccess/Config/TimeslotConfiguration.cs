using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using scheduler_service.DataAccess.Entities;

namespace scheduler_service.DataAccess.Config;

public class TimeslotConfiguration : IEntityTypeConfiguration<Timeslot>
{
    public void Configure(EntityTypeBuilder<Timeslot> builder)
    {
        builder.ToTable("timeslot");

        builder.Property(e => e.Id)
        .HasColumnName("id");

        builder.HasKey(e => e.Id);

        builder.Property(e => e.ApartmentId)
        .HasColumnName("apartment_id")
        .IsRequired();

        builder.Property(e => e.StartDatetime)
        .HasColumnName("startdatetime")
        .IsRequired();

        builder.Property(e => e.EndDatetime)
        .HasColumnName("enddatetime")
        .IsRequired();

        builder.HasIndex(e => e.ApartmentId);

        builder.HasIndex(e => new { e.StartDatetime, e.EndDatetime });

        builder.HasMany(e => e.Reservations)
        .WithOne(e => e.Timeslot);
    }
}